package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/schraf/assistant/pkg/models"
	"google.golang.org/api/option"
)

func startJob(ctx context.Context, contentType string, request models.ContentRequest) error {
	jobName := os.Getenv("CLOUD_RUN_JOB_NAME")
	if jobName == "" {
		return fmt.Errorf("CLOUD_RUN_JOB_NAME environment variable is not set")
	}

	region := os.Getenv("CLOUD_RUN_JOB_REGION")
	if region == "" {
		return fmt.Errorf("CLOUD_RUN_JOB_REGION environment variable is not set")
	}

	projectId := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectId == "" {
		return fmt.Errorf("GOOGLE_CLOUD_PROJECT environment variable is not set")
	}

	//--==================================================================--
	//--== ENCODE THE REQUEST BODY
	//--==================================================================--

	requestBodyJson, err := json.Marshal(request.Body)
	if err != nil {
		return fmt.Errorf("failed to marshal content request: %w", err)
	}

	encodedRequestBody := base64.StdEncoding.EncodeToString(requestBodyJson)

	//--==================================================================--
	//--== REQUEST THE CLOUD RUN JOB
	//--==================================================================--

	client, err := run.NewJobsClient(ctx, option.WithQuotaProject(projectId))
	if err != nil {
		return fmt.Errorf("failed to create Cloud Run Jobs client: %w", err)
	}
	defer client.Close()

	// Build the job execution request
	jobPath := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectId, region, jobName)

	// Execute the job with the request in an environment variable
	jobRequest := &runpb.RunJobRequest{
		Name: jobPath,
		Overrides: &runpb.RunJobRequest_Overrides{
			ContainerOverrides: []*runpb.RunJobRequest_Overrides_ContainerOverride{
				{
					Env: []*runpb.EnvVar{
						{
							Name: "CONTENT_TYPE",
							Values: &runpb.EnvVar_Value{
								Value: contentType,
							},
						},
						{
							Name: "REQUEST_ID",
							Values: &runpb.EnvVar_Value{
								Value: request.Id.String(),
							},
						},
						{
							Name: "REQUEST_BODY",
							Values: &runpb.EnvVar_Value{
								Value: encodedRequestBody,
							},
						},
					},
				},
			},
		},
	}

	_, err = client.RunJob(ctx, jobRequest)
	if err != nil {
		return fmt.Errorf("failed to execute Cloud Run Job: %w", err)
	}

	return nil
}
