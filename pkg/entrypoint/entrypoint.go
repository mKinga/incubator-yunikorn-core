/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package entrypoint

import (
	"github.com/apache/incubator-yunikorn-core/pkg/cache"
	"github.com/apache/incubator-yunikorn-core/pkg/handler"
	"github.com/apache/incubator-yunikorn-core/pkg/rmproxy"
	"github.com/apache/incubator-yunikorn-core/pkg/scheduler"
	"github.com/apache/incubator-yunikorn-core/pkg/webservice"
)

// options used to control how services are started
type StartupOptions struct {
	manualScheduleFlag bool
	startWebAppFlag    bool
}

func StartAllServices() *ServiceContext {
	return startAllServicesWithParameters(
		StartupOptions{
			manualScheduleFlag: false,
			startWebAppFlag:    true,
		})
}

// Visible by tests
func StartAllServicesWithManualScheduler() *ServiceContext {
	return startAllServicesWithParameters(
		StartupOptions{
			manualScheduleFlag: true,
			startWebAppFlag:    false,
		})
}

func startAllServicesWithParameters(opts StartupOptions) *ServiceContext {
	cache := cache.NewClusterInfo()
	scheduler := scheduler.NewScheduler(cache)
	proxy := rmproxy.NewRMProxy()

	eventHandler := handler.EventHandlers{
		CacheEventHandler:     cache,
		SchedulerEventHandler: scheduler,
		RMProxyEventHandler:   proxy,
	}

	// start services
	cache.StartService(eventHandler)
	scheduler.StartService(eventHandler, opts.manualScheduleFlag)
	proxy.StartService(eventHandler)

	context := &ServiceContext{
		RMProxy:   proxy,
		Cache:     cache,
		Scheduler: scheduler,
	}

	if opts.startWebAppFlag {
		webapp := webservice.NewWebApp(cache, scheduler)
		webapp.StartWebApp()
		context.WebApp = webapp
	}

	return context
}
