// Package state manages persistent snapshots of observed open ports for
// portwatch.
//
// # Overview
//
// The state package decouples port observation from change detection by
// providing a simple file-backed store. Each time the daemon completes a
// scan cycle it saves a Snapshot to disk. On the next cycle the previous
// Snapshot is loaded and compared against the new one using Diff, which
// returns the ports that have been opened or closed in the interval.
//
// # Usage
//
//	store := state.New("/var/lib/portwatch/state.json")
//
//	// Load the last known state (empty on first run).
//	prev, err := store.Load()
//
//	// ... perform a scan to build curr ...
//
//	// Detect changes.
//	opened, closed := state.Diff(prev, curr)
//
//	// Persist the new state.
//	err = store.Save(curr)
//
// # File Format
//
// Snapshots are stored as indented JSON and the file is created with mode
// 0600 to prevent other users from reading potentially sensitive network
// topology information.
package state
