package bootstrap

import "treehole_next/apis/hole"

func startTasks() chan struct{} {
	done := make(chan struct{}, 1)
	go hole.UpdateHoleViews(done)
	return done
}
