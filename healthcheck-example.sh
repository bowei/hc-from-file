#!bin/bash

# Copyright 2022 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Example healthchecking script
#
# Loops checking endpoints for 200. If any endpoints are not 200, will write a
# 500 response to the resultFile.
#
resultFile="/result"

echo "500 not initialized" > "${resultFile}"

while true; do
  healthy=1
  for endpoint in \
    localhost:9000 \
    localhost:9001
  do
    if ! curl --max-time 5 --silent "${endpoint}"; then
      healthy=0
      echo "${endpoint} was not healthy" | ts
    fi
  done

  if [[ "${healthy}" == "1" ]]; then
    echo "200 ok" > "${resultFile}"
    echo "All healthchecks passed" | ts
  else
    echo "Healthchecks failed" | ts
    echo "500 not ok" > "${resultFile}"
  fi

  sleep 1
done
