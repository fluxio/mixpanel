Mixpanel Data Export API Go Client 

License
=======

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.


Using Mixpanel Data Export Go client
====================

Mixpanel package offers a client with a very simple API. Its example usage
can be found in the example.go file.

Basically, one needs to specify Mixpanel account's ApiKey and ApiSecret for
the client to use while making requests.

In order to run an example code type:

    $ go run example.go

It is expected to return a server response containing an "Invalid API key" error.


Go client's API
============================

The most important method is Request(...) which forms a request and sends it to the
Mixpanel's server. It has the following signature:

func (m Mixpanel) Request(
    methods []string,
    params map[string](interface{}),
    http_method string,
    format string,
) ([]byte, error)

, where:
methods     -   slice of strings indicating a particular endpoint. For example,
                []string{"events"} hits https://mixpanel.com/api/2.0/events

params      -   map of parameters required by an endpoint. For example,
                map[string](interface{}){"event":[]string{"pages", "home"}} contains one
                parameter "event" with a slice value of "pages" and "home".

http_method -   specifies which HTTP method should be used. For GET, parameters are
                passed along in the URL, for other methods request Body is used instead.

format      -   The format of returned data, it should be either "json" or "csv".


A note about using Go client
============================

Mixpanel's Data Export API endpoints behave in a peculiar way and are very sensitive
to the format of passed along parameters. For example, even though Request(...) requires
a map of [string](interface{}), Mixpanel's endpoints would not accept arbitrary structures,
say, another map internally although slices are fine.
