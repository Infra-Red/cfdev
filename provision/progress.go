package provision

import (
	"fmt"
	"time"

	"code.cloudfoundry.org/cfdev/bosh"
	"code.cloudfoundry.org/cfdev/errors"
)

func (c *Controller) report(start time.Time, ui UI, b *bosh.Bosh, service Service, errChan chan error) error {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				return errors.SafeWrap(err, fmt.Sprintf("Failed to deploy %s", service.Name))
			}

			ui.Writer().Write([]byte(fmt.Sprintf("\r\033[K  Done (%s)\n", time.Now().Sub(start).Round(time.Second))))
			return nil
		case <-ticker.C:
			p := b.GetVMProgress(start, service.Deployment, service.IsErrand)

			switch p.State {
			case bosh.Preparing:
				ui.Writer().Write([]byte(fmt.Sprintf("\r\033[K  Preparing deployment (%s)", p.Duration.Round(time.Second))))
			case bosh.Deploying:
				ui.Writer().Write([]byte(fmt.Sprintf("\r\033[K  Progress: %d of %d (%s)", p.Done, p.Total, p.Duration.Round(time.Second))))
			case bosh.RunningErrand:
				ui.Writer().Write([]byte(fmt.Sprintf("\r\033[K  Running errand (%s)", p.Duration.Round(time.Second))))
			}
		}
	}
}
