package main

import (
	"encoding/json"
	_ "errors"
	_ "fmt"

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
						ID       string `json:"id"`
						Title    string `json:"title"`
						DateTime string `json:"dateTime"`
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
			// ID: edge.Node.ID,
			// EventDate: edge.Node.DateTime,
			EventLocation:    edge.Node.Venue.Address,
			EventDescription: "?",
			EventGroupName:   groupUrlName,
			GroupName:        groupUrlName,
			Dynamic:          true,
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
//
// // VerifyGroupParameters verifies that the group is in Calgary and has Technology as its topic
// func (m *MeetupService) VerifyGroupParameters(groupName string, accessToken string) error {
// 	// Create GraphQL query
// 	graphqlReq := GraphQLRequest{
// 		Query: `
//             query GetEventsByGroup($groupUrlname: String!) {
//                 groupByUrlname(urlname: $groupUrlname) {
//                     id
//                     name
//                     city
//                     topicCategory {
//                         id
//                         urlkey
//                         name
//                         color
//                         imageUrl
//                         defaultTopic {
//                             name
//                         }
//                     }
//                 }
//             }
//         `,
// 		Variables: map[string]interface{}{
// 			"groupUrlname": groupName,
// 		},
// 	}
//
// 	// Marshal request to JSON
// 	graphqlReqBytes, err := json.Marshal(graphqlReq)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Create an HTTP client using Fiber's Agent
// 	agent := fiber.AcquireAgent()
// 	defer fiber.ReleaseAgent(agent)
//
// 	// Setup the request
// 	req := agent.Request()
// 	req.Header.Set("Authorization", "Bearer "+accessToken)
// 	req.Header.Set("Content-Type", "application/json")
// 	req.SetRequestURI(MeetupGraphqlApiUrl)
// 	req.SetBody(graphqlReqBytes)
// 	req.Header.SetMethod(fiber.MethodPost)
//
// 	// Send the request
// 	if err := agent.Parse(); err != nil {
// 		return err
// 	}
//
// 	// Get the response
// 	resp, errs := agent.Bytes()
// 	if len(errs) > 0 {
// 		return errs[0]
// 	}
//
// 	// Parse response
// 	var result map[string]interface{}
// 	if err := json.Unmarshal(resp, &result); err != nil {
// 		return err
// 	}
//
// 	// Check for errors in the GraphQL response
// 	if errNodes, exists := result["errors"].([]interface{}); exists && len(errNodes) > 0 {
// 		if errMsg, ok := errNodes[0].(map[string]interface{})["message"].(string); ok {
// 			return GraphQLError{Message: errMsg}
// 		}
// 		return GraphQLError{Message: "Unknown GraphQL error"}
// 	}
//
// 	// Extract city and topic category
// 	data, ok := result["data"].(map[string]interface{})
// 	if !ok {
// 		return errors.New("invalid response format: missing data field")
// 	}
//
// 	groupByUrlname, ok := data["groupByUrlname"].(map[string]interface{})
// 	if !ok {
// 		return errors.New("invalid response format: missing groupByUrlname field")
// 	}
//
// 	// Check city
// 	city, ok := groupByUrlname["city"].(string)
// 	if !ok {
// 		return errors.New("invalid response format: missing city field")
// 	}
//
// 	if city != "Calgary" {
// 		return InvalidGroupParameterError{Message: fmt.Sprintf("Group from wrong city provided. Provided '%s'", city)}
// 	}
//
// 	// Check topic category
// 	topicCategory, ok := groupByUrlname["topicCategory"].(map[string]interface{})
// 	if !ok {
// 		return errors.New("invalid response format: missing topicCategory field")
// 	}
//
// 	topicCategoryName, ok := topicCategory["name"].(string)
// 	if !ok {
// 		return errors.New("invalid response format: missing topicCategory.name field")
// 	}
//
// 	if topicCategoryName != "Technology" {
// 		return InvalidGroupParameterError{Message: fmt.Sprintf("Group from wrong topic category provided. Provided '%s'", topicCategoryName)}
// 	}
//
// 	return nil
// }
