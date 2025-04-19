package main

import (
	"encoding/json"
	_ "errors"
	_ "fmt"
	"time"

	"github.com/gofiber/fiber/v3/client"
)

const MeetupGraphqlApiUrl = "https://api.meetup.com/gql"

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type GraphQLError struct {
	Message string
}

type InvalidGroupParameterError struct {
	Message string
}

func (e InvalidGroupParameterError) Error() string {
	return e.Message
}

type MeetupResponse struct {
	Data struct {
		GroupByUrlname struct {
			UpcomingEvents struct {
				Edges []struct {
					Node struct {
						ID       string    `json:"id"`
						Title    string    `json:"title"`
						DateTime time.Time `json:"dateTime"`
						Venue    struct {
							Name    string `json:"name"`
							Address string `json:"address"`
						} `json:"venue"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"upcomingEvents"`
			City          string `json:"city"`
			TopicCategory struct {
				Name string `json:"name"`
			} `json:"topicCategory"`
		} `json:"groupByUrlname"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

func FetchEvents(accessToken string, groupUrlName string) ([]Event, error) {
	graphqlReq := GraphQLRequest{
		Query: `
			query ($groupUrlname: String!, $input: ConnectionInput!) { 
				groupByUrlname(urlname: $groupUrlname) { 
					upcomingEvents(input: $input) { 
						edges { 
							node { 
								id 
								title 
								dateTime 
								venue { 
									name 
									address 
								} 
							} 
						} 
					} 
				} 
			}
		`,
		Variables: map[string]any{
			"groupUrlname": groupUrlName,
			"input": map[string]any{
				"first": 10,
			},
		},
	}

	graphqlReqBytes, err := json.Marshal(graphqlReq)
	if err != nil {
		return nil, err
	}

	cc := client.New()
	rsp, err := cc.Post(MeetupGraphqlApiUrl, client.Config{
		Header: map[string]string{
			"Authorization": "Bearer " + accessToken,
			"Content-Type":  "application/json",
		},
		Body: graphqlReqBytes,
	})
	if err != nil {
		return nil, err
	}

	var decodedRsp MeetupResponse
	if err := json.Unmarshal(rsp.Body(), &decodedRsp); err != nil {
		return nil, err
	}

	eventsToSave := make([]Event, 0, len(decodedRsp.Data.GroupByUrlname.UpcomingEvents.Edges))
	for _, edge := range decodedRsp.Data.GroupByUrlname.UpcomingEvents.Edges {
		event := Event{
			Date:        edge.Node.DateTime,
			Location:    edge.Node.Venue.Address,
			Description: "?",
			// EventGroupName:   groupUrlName,
			// GroupName:        groupUrlName,
			// Dynamic:          true,
		}

		eventsToSave = append(eventsToSave, event)
	}

	return eventsToSave, nil
}

// func (m *MeetupService) SyncEventsForGroup(accessToken string, groupName string) error {
// 	oldEvents, err := m.eventService.GetDynamicEventsForGroup(groupName)
// 	if err != nil {
// 		return err
//
//
// 	newEvents, err := m.FetchEvents(accessToken, groupName)
// 	if err != nil {
// 		return err
// 	}
//
// 	oldEventsMap := make(map[string]Event)
// 	for _, event := range oldEvents {
// 		oldEventsMap[event.EventID] = event
// 	}
//
// 	for i := range newEvents {
// 		if oldEvent, exists := oldEventsMap[newEvents[i].EventID]; exists {
// 			newEvents[i].ID = oldEvent.ID
// 		}
// 	}
//
// 	if err := m.eventService.SaveEvents(newEvents); err != nil {
// 		return err
// 	}
//
// 	eventsToRemove := findEventsToRemove(oldEvents, newEvents)
//
// 	if len(eventsToRemove) > 0 {
// 		return m.eventService.DeleteEvents(eventsToRemove)
// 	}
//
// 	return nil
// }

type GroupParametersResponse struct {
	Data struct {
		GroupByUrlname struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			City          string `json:"city"`
			TopicCategory struct {
				ID           string `json:"id"`
				URLKey       string `json:"urlkey"`
				Name         string `json:"name"`
				Color        string `json:"color"`
				ImageURL     string `json:"imageUrl"`
				DefaultTopic struct {
					Name string `json:"name"`
				} `json:"defaultTopic"`
			} `json:"topicCategory"`
		} `json:"groupByUrlname"`
	} `json:"data"`
}

type GroupParameters struct {
	City  string
	Topic string
}

// VerifyGroupParameters verifies that the group is in Calgary and has Technology as its topic
func GetGroupParameters(groupName string, accessToken string) (GroupParameters, error) {
	// Create GraphQL query
	graphqlReq := GraphQLRequest{
		Query: `
            query GetEventsByGroup($groupUrlname: String!) {
                groupByUrlname(urlname: $groupUrlname) {
                    id
                    name
                    city
                    topicCategory {
                        id
                        urlkey
                        name
                        color
                        imageUrl
                        defaultTopic {
                            name
                        }
                    }
                }
            }
        `,
		Variables: map[string]interface{}{
			"groupUrlname": groupName,
		},
	}

	// Marshal request to JSON
	graphqlReqBytes, err := json.Marshal(graphqlReq)
	if err != nil {
		return GroupParameters{}, err
	}

	cc := client.New()
	rsp, err := cc.Post(MeetupGraphqlApiUrl, client.Config{
		Header: map[string]string{
			"Authorization": "Bearer " + accessToken,
			"Content-Type":  "application/json",
		},
		Body: graphqlReqBytes,
	})
	if err != nil {
		return GroupParameters{}, err
	}

	var parsedRsp GroupParametersResponse
	if err := json.Unmarshal(rsp.Body(), &parsedRsp); err != nil {
		return GroupParameters{}, err
	}
	return GroupParameters{
		City:  parsedRsp.Data.GroupByUrlname.City,
		Topic: parsedRsp.Data.GroupByUrlname.TopicCategory.Name,
	}, nil
}

	// todo: keep
	// if city != "Calgary" {
	// 	return InvalidGroupParameterError{Message: fmt.Sprintf("Group from wrong city provided. Provided '%s'", city)}
	// }
	// if topicCategoryName != "Technology" {
	// 	return InvalidGroupParameterError{Message: fmt.Sprintf("Group from wrong topic category provided. Provided '%s'", topicCategoryName)}
	// }
