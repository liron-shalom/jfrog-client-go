package services

import (
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type ReplicationService struct {
	isUpdate   bool
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewReplicationService(client *rthttpclient.ArtifactoryHttpClient, isUpdate bool) *ReplicationService {
	return &ReplicationService{client: client, isUpdate: isUpdate}
}

func (rs *ReplicationService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return rs.client
}

func (rs *ReplicationService) PerformRequest(replicationParams []byte, repoKey string) error {
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.replications.ReplicationConfigRequest+json", &httpClientsDetails.Headers)
	var url = rs.ArtDetails.GetUrl() + "api/replications/" + repoKey
	var operationString string
	var resp *http.Response
	var body []byte
	var err error
	if rs.isUpdate {
		log.Info("Update replication job...")
		operationString = "updating"
		resp, body, err = rs.client.SendPost(url, replicationParams, &httpClientsDetails)
	} else {
		log.Info("Creating replication job...")
		operationString = "creating"
		resp, body, err = rs.client.SendPut(url, replicationParams, &httpClientsDetails)
	}
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done " + operationString + " repository.")
	return nil
}

type CommonReplicationParams struct {
	CronExp                string `json:"cronExp,omitempty"`
	RepoKey                string `json:"repoKey,omitempty"`
	EnableEventReplication bool   `json:"enableEventReplication,omitempty"`
	SocketTimeoutMillis    int    `json:"socketTimeoutMillis,omitempty"`
	Enabled                bool   `json:"enabled,omitempty"`
	SyncDeletes            bool   `json:"syncDeletes,omitempty"`
	SyncProperties         bool   `json:"syncProperties,omitempty"`
	SyncStatistics         bool   `json:"syncStatistics,omitempty"`
	PathPrefix             string `json:"pathPrefix,omitempty"`
}

type PullReplicationParams struct {
	CommonReplicationParams
}

type PushReplicationParams struct {
	CommonReplicationParams
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	URL      string `json:"url,omitempty"`
}
