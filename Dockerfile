ARG AGENT_VERSION=9.0.1

FROM docker.elastic.co/elastic-agent/elastic-agent:${AGENT_VERSION}

COPY build/elastic-agent /usr/share/elastic-agent/elastic-agent
COPY build/elastic-agent.yml /usr/share/elastic-agent/elastic-agent.yml
COPY .build_hash.txt /usr/share/elastic-agent/.build_hash.txt.new
RUN mv /usr/share/elastic-agent/data/elastic-agent-$(cat /usr/share/elastic-agent/.build_hash.txt| cut -c 1-6) /usr/share/elastic-agent/data/elastic-agent-$(cat /usr/share/elastic-agent/.build_hash.txt.new| cut -c 1-6) && \
  ln -s -f /usr/share/elastic-agent/data/elastic-agent-$(cat /usr/share/elastic-agent/.build_hash.txt.new| cut -c 1-6)/elastic-agent /usr/share/elastic-agent/elastic-agent && \
  mv /usr/share/elastic-agent/.build_hash.txt /usr/share/elastic-agent/.build_hash.txt.old && \
  mv /usr/share/elastic-agent/.build_hash.txt.new /usr/share/elastic-agent/.build_hash.txt
