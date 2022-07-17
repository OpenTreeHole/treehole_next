package bootstrap

import "treehole_next/apis/hole"

func startTasks() {
	go hole.UpdateHoleViews()
}
