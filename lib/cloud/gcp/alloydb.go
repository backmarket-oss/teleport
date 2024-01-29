/*
 * Teleport
 * Copyright (C) 2023  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package gcp

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	alloydb "cloud.google.com/go/alloydb/apiv1/alloydbpb"
	"github.com/gravitational/trace"

	"github.com/gravitational/teleport/api/constants"
	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/tlsca"
)

// AlloyDBAdminClient defines an interface providing access to the GCP AlloyDB API.
type AlloyDBAdminClient interface {
	// UpdateUser updates an existing user for the project/instance configured in a session.
	UpdateUser(ctx context.Context, db types.Database, dbUser string, user *alloydb.User) error
	// GetCluster returns database cluster details for the project/instance
	// configured in a session.
	GetCluster(ctx context.Context, db types.Database) (*alloydb.Instance, error)
	// GenerateEphemeralCert returns a new client certificate with RSA key for the
	// project/instance configured in a session.
	GenerateEphemeralCert(ctx context.Context, db types.Database, identity tlsca.Identity) (*tls.Certificate, error)
}

// NewAlloyDBClient returns an AlloyDBClient interface wrapping sqladmin.Service.
func NewAlloyDBClient(ctx context.Context) (AlloyDBAdminClient, error) {
	service, err := alloydb.NewAlloyDBAdminClient()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return &gcpAlloyDBAdminClient{service: service}, nil
}

// gcpAlloyDBAdminClient implements the GCPAlloyDBAdminClient interface by wrapping
// sqladmin.Service.
type gcpAlloyDBAdminClient struct {
	service *alloydb.AlloyDBAdminClient
}

// UpdateUser updates an existing user in an AlloyDB for the project/instance
// configured in a session.
func (g *gcpAlloyDBClient) UpdateUser(ctx context.Context, db types.Database, dbUser string, user *alloydb.User) error {
	clusterUserService := alloydb.NewProjectsLocationsClustersUsersService(g.service)
	_, err := clusterUserService.Patch(
		db.GetGCP().ProjectID,
		db.GetGCP().InstanceID,
		user).Name(dbUser).Host("%").Context(ctx).Do()
	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

// GetCluster returns database instance details from Cloud SQL for the
// project/instance configured in a session.
func (g *gcpAlloyDBClient) GetCluster(ctx context.Context, db types.Database) (*alloydb.Instance, error) {
	gcp := db.GetGCP()
	dbi, err := g.service.Cluster.Get(gcp.ProjectID, gcp.InstanceID).Context(ctx).Do()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return dbi, nil
}

// GenerateEphemeralCert returns a new client certificate with RSA key created
// using the GenerateEphemeralCertRequest Cloud SQL API. Client certificates are
// required when enabling SSL in Cloud SQL.
func (g *gcpAlloyDBClient) GenerateEphemeralCert(ctx context.Context, db types.Database, identity tlsca.Identity) (*tls.Certificate, error) {
	// TODO(jimbishopp): cache database certificates to avoid expensive generate
	// operation on each connection.

	// Generate RSA private key, x509 encoded public key, and append to certificate request.
	pkey, err := rsa.GenerateKey(rand.Reader, constants.RSAKeySize)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	pkix, err := x509.MarshalPKIXPublicKey(pkey.Public())
	if err != nil {
		return nil, trace.Wrap(err)
	}

	// Make API call.
	gcp := db.GetGCP()
	req := g.service.Connect.GenerateEphemeralCert(gcp.ProjectID, gcp.InstanceID, &alloydb.GenerateClientCertificateRequest{
		PublicKey:    string(pem.EncodeToMemory(&pem.Block{Bytes: pkix, Type: "RSA PUBLIC KEY"})),
		CertDuration: fmt.Sprintf("%ds", int(time.Until(identity.Expires).Seconds())),
	})
	resp, err := req.Context(ctx).Do()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	// Create TLS certificate from returned ephemeral certificate and private key.
	cert, err := tls.X509KeyPair([]byte(resp.EphemeralCert.Cert), tlsca.MarshalPrivateKeyPEM(pkey))
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return &cert, nil
}
