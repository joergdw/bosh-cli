package templatescompiler

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bmproperty "github.com/cloudfoundry/bosh-micro-cli/common/property"
	bmreljob "github.com/cloudfoundry/bosh-micro-cli/release/job"
)

type JobListRenderer interface {
	Render(
		releaseJobs []bmreljob.Job,
		jobProperties bmproperty.Map,
		globalProperties bmproperty.Map,
		deploymentName string,
	) (RenderedJobList, error)
}

type jobListRenderer struct {
	jobRenderer JobRenderer
	logger      boshlog.Logger
	logTag      string
}

func NewJobListRenderer(
	jobRenderer JobRenderer,
	logger boshlog.Logger,
) JobListRenderer {
	return &jobListRenderer{
		jobRenderer: jobRenderer,
		logger:      logger,
		logTag:      "jobListRenderer",
	}
}

func (r *jobListRenderer) Render(
	releaseJobs []bmreljob.Job,
	jobProperties bmproperty.Map,
	globalProperties bmproperty.Map,
	deploymentName string,
) (RenderedJobList, error) {
	r.logger.Debug(r.logTag, "Rendering job list: deploymentName='%s' jobProperties=%#v globalProperties=%#v", deploymentName, jobProperties, globalProperties)
	renderedJobList := NewRenderedJobList()

	// render all the jobs' templates
	for _, releaseJob := range releaseJobs {
		renderedJob, err := r.jobRenderer.Render(releaseJob, jobProperties, globalProperties, deploymentName)
		if err != nil {
			defer renderedJobList.DeleteSilently()
			return renderedJobList, bosherr.WrapErrorf(err, "Rendering templates for job %#v", releaseJob)
		}
		renderedJobList.Add(renderedJob)
	}

	return renderedJobList, nil
}
