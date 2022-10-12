package impl

import (
	"context"
	"github.com/zginkgo/ginkgo_cmdb/apps/secret"
)

type secretMock struct {
	secret.UnimplementedServiceServer
}

func (m *secretMock) CreateSecret(ctx context.Context, request *secret.CreateSecretRequest) (*secret.Secret, error) {
	return nil, nil
}

func (m *secretMock) QuerySecret(context.Context, *secret.QuerySecretRequest) (*secret.SecretSet, error) {
	return nil, nil
}

func (m *secretMock) DescribeSecret(context.Context, *secret.DescribeSecretRequest) (
	*secret.Secret, error) {
	ins := secret.NewDefaultSecret()
	ins.Data.ApiKey = "AKIDTilY5LJWUCbIlyHP2PVK3WiyBJBoKoDL"
	ins.Data.ApiSecret = "zhZCZgW8WdTYP0LP6g2DDXcdOvDo8C3f"
	return ins, nil
}

func (m *secretMock) DeleteSecret(context.Context, *secret.DeleteSecretRequest) (*secret.Secret, error) {
	return nil, nil
}
