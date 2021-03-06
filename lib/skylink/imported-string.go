package skylink

import (
	"errors"
	"fmt"
	"log"

	"github.com/stardustapp/dustgo/lib/base"
	"github.com/stardustapp/dustgo/lib/inmem"
)

type importedString struct {
	svc  *nsimport
	path string
	*inmem.String
}

var _ base.String = (*importedString)(nil)

// Support subscription passthrough
var _ Subscribable = (*importedString)(nil)

// TODO: DRY this the fuck up
func (e *importedString) Subscribe(s *Subscription) (err error) {
	req := nsRequest{
		Op:    "subscribe",
		Path:  e.path,
		Depth: s.MaxDepth,
	}

	resp, err := e.svc.transport.exec(req)
	if err != nil {
		return errors.New("nsimport string subscribe err:")
	}

	if resp.Channel != nil {
		go func(inC <-chan nsResponse, outC chan<- Notification, stopC <-chan struct{}) {
			log.Println("imported-string: Starting subscription pump from", e.path)
		feedLoop:
			for {
				select {
				case pkt, ok := <-inC:
					if !ok {
						log.Println("imported-folder: Propagating sub close downstream")
						break feedLoop
					}

					if pkt.Output != nil && pkt.Output.Name == "notif" {
						var notif Notification
						for _, field := range pkt.Output.Children {
							switch field.Name {
							case "type":
								notif.Type = field.StringValue
							case "path":
								notif.Path = field.StringValue
							case "entry":
								notif.Entry = convertEntryFromApi(&field)
							default:
								log.Println("imported-string WARN: sub got weird Next field,", field.Name)
							}
						}
						//log.Println("imported-string: sub notification:", notif)
						outC <- notif
					}

					// the above assumes that the remote can't double-terminate,
					// and will close immediately after any terminal event

				case <-stopC:
					log.Println("imported-string: Propagating sub stop upstream")
					resp, err := e.svc.transport.exec(nsRequest{
						Op:   "stop",
						Path: fmt.Sprintf("/chan/%d", resp.Chan),
					})
					if err != nil {
						log.Println("WARN: nsimport folder stop err:", err)
					} else {
						log.Println("nsimport folder stop happened:", resp)
					}
					// stop checking the stop chan, just finish out the main chan
					stopC = nil

				}
			}
			log.Println("imported-string: Completed subscription pump from", e.path)
			close(outC)
		}(resp.Channel, s.StreamC, s.StopC)
	}

	if resp.Status == "Ok" {
		return nil
	} else if resp.Output != nil && resp.Output.Type == "String" {
		return errors.New("Subscription failed. Cause: " + resp.Output.StringValue)
	} else {
		return errors.New("Subscription attempt wasn't Ok, was " + resp.Status)
	}
}
