# not used by Tilt (yet)
FROM python:3.10-slim

WORKDIR /app

# RUN apk add --no-cache gcc musl-dev python3-dev libffi-dev openssl-dev cargo g++

ENV PYTHONFAULTHANDLER=1 \
  PYTHONUNBUFFERED=1 \
  PYTHONHASHSEED=random \
  PIP_NO_CACHE_DIR=off \
  PIP_DISABLE_PIP_VERSION_CHECK=on \
  PIP_DEFAULT_TIMEOUT=100 \
  POETRY_VERSION=1.3.2

RUN pip install "poetry==$POETRY_VERSION"
# This image is single purpose, so we won't need to compartmentalize Py deps in virtualenvs.
RUN poetry config virtualenvs.create false

COPY pyproject.toml pyproject.toml
RUN poetry install
COPY . .

EXPOSE 5678
CMD ["opentelemetry-instrument", "flask", "run", "--host", "0.0.0.0", "--port", "5678"]