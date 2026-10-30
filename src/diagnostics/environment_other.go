//go:build !linux

package diagnostics

func platformEnvironment(env *Environment) {}

func describeVolume(volume *Volume) {}
