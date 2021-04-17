package cronjob

import (
	"time"

	"github.com/rs/zerolog/log"
)

// Run executes a cronjob
func Run() {
	timeSince := time.Now().Add(-15 * time.Minute)
	log.Info().Msgf("cronjob running %s.", timeSince)
}
