// Package gococ is a Chain-of-Command (CoCo) framework for writing
// applications.
//
// # Tutorials
//
// * Building Web Apps with Gococ: https://github.com/Khulnasoft-lab/gococ-web-tutorial
//
// * Building CLI Apps with Gococ: https://github.com/Khulnasoft-lab/gococ-cli-tutorial
//
// A chain of command framework works as follows:
//
//   - A "route" is constructed as a chain of commands -- a series of
//     single-purpose tasks that are run in sequence.
//   - An application is composed of one or more routes.
//   - Commands in a route communicate using a Context.
//   - An application Router is used to receive a route name and then
//     execute the appropriate chain of commands.
//
// To create a new Gococ application, use gococ.Gococ(). This will
// configure and create a new registry, request router, and context.
// From there, use the Registry to build chains of commands, and then
// use the Router to execute chains of commands.
//
// Unlike other CoCo implementations (like Pronto.js or Fortissimo),
// Gococ commands are just functions.
//
// # Interrupts
//
// There are four types of interrupts that you may wish to return:
//
//  1. FatalError: This will stop the route immediately.
//  2. RecoverableError: This will allow the route to continue moving.
//  3. Stop: This will stop the current request, but not as an error.
//  4. Reroute: This will stop executing the current route, and switch to executing another route.
//
// To learn how to write Gococ applications, you may wish to examine
// the small Skunk application: https://github.com/technosophos/skunk.
package gococ

// VERSION provides the current version of Gococ.
const VERSION = "1.3.0"

// Gococ creates a new Gococ app.
//
// This is the main progenitor of a Gococ application. Whether a plain
// Gococ app, or a Web or CLI program, this is the function you will use
// to bootstrap.
//
// The `*Registry` is used to declare new routes, where a "route" may be thought
// of as a task composed of a series of steps (commands).
//
// The `*Router` is responsible for the actual execution of a Gococ route. The
// main method used to call a route is `Router.HandleRequest()`.
//
// The `Context` is a container for passing information down a chain of commands.
// Apps may insert "global" information to a context at startup and make it
// available to all commands.
func Gococ() (reg *Registry, router *Router, cxt Context) {
	cxt = NewContext()
	reg = NewRegistry()
	router = NewRouter(reg)
	return
}

// Command executes a command and returns a result.
// A Gococ app has a registry, which has zero or more routes. Each route
// executes a sequence of zero or more commands. A command is of this type.
type Command func(cxt Context, params *Params) (interface{}, Interrupt)

// Interrupt is a generic return for a command.
// Generally, a command should return one of the following in the interrupt slot:
// - A FatalError, which will stop processing.
// - A RecoverableError, which will continue the chain.
// - A Reroute, which will cause a different route to be run.
type Interrupt interface{}

// Creates a new Reroute.
func NewReroute(route string) *Reroute {
	return &Reroute{route}
}

// Reroute is a command can return a Reroute to tell the router to execute a
// different route.
//
// A `Command` may return a `Reroute` to cause Gococ to stop executing the
// current route and jump to another.
//
//	func Forward(c Context, p *Params) (interface{}, Interrupt) {
//		return nil, &Reroute{"anotherRoute"}
//	}
type Reroute struct {
	Route string
}

// RouteTo returns the route to reroute to.
func (rr *Reroute) RouteTo() string {
	return rr.Route
}

// Stop a route, but not as an error condition.
//
// When Gococ encounters a `Stop`, it will not execute any more commands on a
// given route. However, it will not emit an error, either.
type Stop struct{}

// RecoverableError is an error that should not cause the router to stop processing.
//
// When Gococ encounters a `RecoverableError`, it will log the error as a
// warning, but will then continue to execute the next command in the route.
type RecoverableError struct {
	Message string
}

// Error returns the error message.
func (err *RecoverableError) Error() string {
	return err.Message
}

// FatalError is a fatal error, which will stop the router from continuing a route.
//
// When Gococ encounters a `FatalError`, it will log the error and immediately
// stop processing the route.
//
// Note that by default Gococ treats and unhandled `error` as if it were a
// `FatalError`.
type FatalError struct {
	Message string
}

// Error returns the error message.
func (err *FatalError) Error() string {
	return err.Message
}
