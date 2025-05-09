/*
 * Copyright 2023 Caio Matheus Marcatti Calimério
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package grpc_webserver

import (
	"context"

	"github.com/caiomarcatti12/nanogo/pkg/context_manager"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func correlationIdInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		fcm := context_manager.NewSafeContextManager()

		var correlationID string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok && len(md["x-correlation-id"]) > 0 {
			correlationID = md["x-correlation-id"][0]
		} else {
			correlationID = uuid.New().String()
		}

		ctxWithCorrelation := metadata.AppendToOutgoingContext(ctx, "x-correlation-id", correlationID)

		var resp interface{}
		var err error

		fcm.SetValues(
			map[interface{}]interface{}{"x-correlation-id": correlationID},
			func() {
				resp, err = handler(ctxWithCorrelation, req)
			},
		)

		return resp, err
	}
}
