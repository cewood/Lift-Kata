package lift

import (
	"errors"
	"time"
)

// Direction ..
type Direction int

// Directions ..
const (
	Up Direction = iota
	Down
	None
)

// Call ..
type Call struct {
	Floor     int
	Direction Direction
}

// Lift ..
type Lift struct {
	ID        string
	Floor     int
	Requests  []int
	DoorsOpen bool
}

// CloseDoors ...
func (l *Lift) CloseDoors() {
	l.DoorsOpen = false
}

// OpenDoors ...
func (l *Lift) OpenDoors() {
	l.DoorsOpen = true
}

// GetDirection is a helper function to quickly calculate
//  the direction the lift is traveling based on the current
//  floor of the lift and the first request pending
func (l *Lift) GetDirection() Direction {
	if len(l.Requests) == 0 {
		return None
	} else if l.Requests[0] > l.Floor {
		return Up
	}

	return Down
}

// MoveUp moves the lift upwards one floor at a time
func (l *Lift) MoveUp() error {
	if l.DoorsOpen {
		return errors.New("The lift doors are open, can't move")
	}

	l.Floor++

	return nil
}

// MoveDown moves the lift downwards one floor at a time
func (l *Lift) MoveDown() error {
	if l.DoorsOpen {
		return errors.New("The lift doors are open, can't move")
	}

	l.Floor--

	return nil
}

// Move ...
func (l *Lift) Move() error {
	var err error

	// Check which floor to visit first
	if l.GetDirection() == Up {
		// Move upwards
		err = l.MoveUp()
	} else if l.GetDirection() == Down {
		// Move downwards
		err = l.MoveDown()
	}

	return err
}

// GetNextRequest ...
func (l *Lift) GetNextRequest() *int {
	if len(l.Requests) == 0 {
		return nil
	}

	return &l.Requests[0]
}

// FulfillRequest handles removing a request from the queue
//  if the current floor is present in the requests queue.
//  It checks if the doors are open, and thus returns an
//  error so that this can be caught and handled correctly.
func (l *Lift) FulfillRequest() error {
	if !l.DoorsOpen {
		// You can't fulfill a request if the doors are closed
		return errors.New("doors were closed")
	}

	if l.Requests[0] == l.Floor {
		l.Requests = append(l.Requests[:0], l.Requests[1:]...)
	}

	return nil
}

// NewRequest handles adding a new request to the queue,
//  which it maintains in a kind of balanced/pivot list.
//  Whereby depending on the direction of the lift, the
//  items will be ordered in ascending order, until an
//  item is request which is below the current floor. In
//  that case the item will be added at the end of the
//  the list in descending order. Here is an example of
//  what this would look like...
//
//  going up => 4, 5, 6, 2, 1, 0
//  going down => 6, 5, 4, 3, 8, 9
func (l *Lift) NewRequest(req int) {
	var i int

	if req == l.Floor {
		// Check if req is the current floor
		return
	} else if len(l.Requests) == 0 {
		// Check if requests is empty, i.e. append
		i = 0
	} else {
		for index, value := range l.Requests {
			if l.Requests[0] > l.Floor {
				if req <= value && req > l.Floor && value > l.Floor {
					i = index
					break
				} else if req > value && value < l.Floor {
					i = index
					break
				}
			} else if l.Requests[0] < l.Floor {
				if req >= value && req < l.Floor && value < l.Floor {
					i = index
					break
				} else if req < value && value > l.Floor {
					i = index
					break
				}
			}

			if index+1 == len(l.Requests) {
				// End of requests, i.e. must append
				i = len(l.Requests)
				break
			}
		}
	}

	if i == len(l.Requests) {
		// Appending
		l.Requests = append(l.Requests, req)
	} else if l.Requests[i] != req {
		// Inserting
		l.Requests = append(l.Requests, 0)
		copy(l.Requests[i+1:], l.Requests[i:])
		l.Requests[i] = req
	}
}

// Tick ...
func (l *Lift) Tick() {
	if l.DoorsOpen && l.GetNextRequest() != nil {
		// If the doors are open, then we have already
		//  visited a requested floor. Thus if there is
		//  any remaining requests, we should close the
		//  doors, because on the next tick we'll want
		//  to move to the next requested floor
		l.CloseDoors()
	} else if !l.DoorsOpen && l.GetNextRequest() != nil && *l.GetNextRequest() == l.Floor {
		// Open the doors requested floor
		l.OpenDoors()

		// Fulfill the request
		l.FulfillRequest()
	} else if !l.DoorsOpen && l.GetNextRequest() != nil {
		// If the doors are closed, there are requests,
		//  and the next request isn't for this floor,
		//  then we should move to that floor
		l.Move()
	}
}

// System ..
type System struct {
	floors []int
	lifts  []Lift
	calls  []Call
}

// NewSystem ..
func NewSystem() *System {
	return &System{floors: []int{}, lifts: []Lift{}, calls: []Call{}}
}

// AddFloors ..
func (s *System) AddFloors(floors ...int) {
	s.floors = append(s.floors, floors...)
}

// AddLifts ..
func (s *System) AddLifts(lifts ...Lift) {
	s.lifts = append(s.lifts, lifts...)
}

// AddCalls ..
func (s *System) AddCalls(calls ...Call) {
	s.calls = append(s.calls, calls...)
}

// CallsFor ..
func (s System) CallsFor(floor int) (calls []Call) {
	calls = []Call{}
	for _, c := range s.calls {
		if c.Floor == floor {
			calls = append(calls, c)
		}
	}
	return calls
}

// TickLifts ...
func (s *System) TickLifts() {
	calls := make([]Call, len(s.calls))
	copy(calls, s.calls)

	// Iterate over our lifts
	for _, lift := range s.lifts {
		// Iterate over the calls
		for i, call := range calls {
			// Check if the
			if call.Direction == lift.GetDirection() && call.Floor > lift.Floor && call.Floor < lift.Requests[0] {
				lift.NewRequest(call.Floor)
				calls = append(calls[:i], calls[i+1:]...)
				break
			} else if call.Direction == lift.GetDirection() && call.Floor < lift.Floor && call.Floor > lift.Requests[0] {
				lift.NewRequest(call.Floor)
				calls = append(calls[:i], calls[i+1:]...)
				break
			}
		}

		lift.Tick()
	}
}

// SatisfyCalls ...
func (s *System) SatisfyCalls() {
	for i, call := range s.calls {
		for _, lift := range s.lifts {
			if call.Direction == lift.GetDirection() && call.Floor == lift.Floor && lift.DoorsOpen {
				s.calls = append(s.calls[:i], s.calls[i+1:]...)
				break
			}
		}
	}
}

// Tick ..
func (s System) Tick(count int, duration time.Duration) {
	for {
		// Check for any satisfied calls
		s.SatisfyCalls()

		// Tick lifts
		s.TickLifts()

		// Sleep
		time.Sleep(time.Duration(count) * duration)
	}
}
