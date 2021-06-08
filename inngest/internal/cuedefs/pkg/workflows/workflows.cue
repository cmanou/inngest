package workflows

import (
	"inngest.com/edges"
)

#EdgeMetadata: edges.#Edge

// a workflow is an entire workflow for our app
#Workflow: {
	name:     string

	// The triggers which start a workflow.
	//
	// If this is a scheduled trigger, only one trigger may exist.
	// Workflows triggered by events may contain multiple event triggers which are exclusive -
	// any of these triggers will start a workflow.
	triggers?: [ ...#Trigger]
	actions?: [ ...#Action]
	edges?: [ ...#Edge]
}

// trigger represents the event that starts our the workflow
#Trigger: #EventTrigger | #ScheduleTrigger

#EventTrigger: {
	event: string
}

#ScheduleTrigger: {
	cron: string
}

#Action: {
	// clientID represents the ID of the action as represented by edges and
	// by the frontend's rendering.
	clientID: uint
	clientID: >=1
	// name of the action
	name: string
	// dsn of the action.  eg "com.datosapp.logic.if" to test a predicate
	// or "com.datosapp.comms.email" to send an email
	dsn: string
	// version of the action DSN to run.  If this is undefined, this defaults
	// to the latest version of the action at time of workflow creation.
	version?: uint
	// Metadata about how the action will be used.  Each action requires custom
	// input to work, eg. what data to transform, what email template to use, etc.
	metadata?: [string]: _
}

#Edge: {
	outgoing: uint | "trigger" // Either the action ID or 'trigger'
	incoming: uint
	metadata?: #EdgeMetadata
}