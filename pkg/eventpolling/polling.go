package eventpolling

import (
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/client"
	"github.com/PRYVT/utils/pkg/store/repository"
	"github.com/rs/zerolog/log"
)

type EventPolling struct {
	client       *client.EventSourcingHttpClient
	eventRepo    *repository.EventRepository
	eventHandler EventHanlder
}

func NewEventPolling(client *client.EventSourcingHttpClient, eventRepo *repository.EventRepository, eventHandler EventHanlder) *EventPolling {
	if client == nil || eventRepo == nil || eventHandler == nil {
		return nil
	}
	return &EventPolling{client: client, eventRepo: eventRepo, eventHandler: eventHandler}
}

func (ep *EventPolling) PollEvents() {

	for {
		time.Sleep(100 * time.Millisecond)
		ep.PollEventsUntilEmpty()
	}

}

func (ep *EventPolling) PollEventsUntilEmpty() {
	for {
		eId, err := ep.eventRepo.GetLastEvent()
		if err != nil {
			log.Err(err).Msg("Error while getting last events")
			continue
		}
		events, err := ep.client.GetEventsSince(eId, 2)
		if err != nil {
			log.Err(err).Msg("Error while polling events")
			continue
		}

		for _, event := range events {

			err := ep.eventHandler.HandleEvent(event)
			if err != nil {
				log.Err(err).Msg("Error while processing event")
				break
			}
		}
		if len(events) == 0 {
			return
		}
		//will this break the db consistency if there are going to be multiple instances of this service?
		// probably but if we dont a volume (that both instances use as a db file) this should be fine
		err = ep.eventRepo.ReplaceEvent(events[len(events)-1].Id)
		if err != nil {
			log.Err(err).Msg("Error while replacing event")

		}
	}
}
