# Full Stack Applications - Rule Engine & Weather Monitor

This repository contains two full-stack applications: a Rule Engine and a Weather Monitoring system. Both applications are containerized using Docker for easy deployment and testing.

## Deploy link : 
- Rule Engine with AST: https://rule-engine-sand.vercel.app/
- Weather Monitor : https://weather-monitor-tau.vercel.app/

# Table of Contents

## Overview
- [Full Stack Applications - Rule Engine & Weather Monitor](#full-stack-applications---rule-engine--weather-monitor)
- [Technology Stack](#technology-stack)
  - [Rule Engine Application](#rule-engine-application)
  - [Weather Monitor Application](#weather-monitor-application)

## Features
- [Application Features](#application-features)
  - [Rule Engine Application](#rule-engine-application-1)
    - [Core Features](#core-features)
    - [Query Optimization](#query-optimization)
    - [Advanced Rule Combination Heuristics](#advanced-rule-combination-heuristics)
  - [Weather Monitoring Application](#weather-monitoring-application)
    - [Core Features](#core-features-1)
    - [Data Analysis](#data-analysis)
    - [Advanced Features](#advanced-features)
    - [Performance Optimization](#performance-optimization)

## Setup & Installation
- [Getting Started](#getting-started)
  - [Option 1: Using Docker Compose](#option-1-using-docker-compose)
  - [Option 2: Using Pre-built Docker Image](#option-2-using-pre-built-docker-image)

## Technical Details
- [Application Ports](#application-ports)
  - [Rule Engine](#rule-engine)
  - [Weather Monitor](#weather-monitor)


## Main  : 

## Technology Stack

### Rule Engine Application
- **Frontend**: Next.js
- **Backend**: Golang (Gin Framework)
- **Database**: MongoDB

### Weather Monitor Application
- **Frontend**: Next.js
- **Backend**: Golang (Gin Framework)
- **Database**: PostgreSQL (TimescaleDB)


## Application Features

### Rule Engine Application
#### Core Features:
- Rule Creation and Management
  - Create individual rules with multiple conditions
  - Combine 2+ rules with complex logical operations
  - Real-time rule evaluation based on input data

#### Query Optimization:
- Optimized queries for:
  - Rule creation
  - Rule combination
  - Rule evaluation
  - Rule optimization

#### Advanced Rule Combination Heuristics:
- Intelligent rule merging strategy:
  - Identifies and isolates common conditions for single evaluation
  - Preserves rule semantics during combination
  - Rules without common conditions: Combined using AND operations
  - Rules with common conditions: Merged using OR operations after extracting common conditions

### Weather Monitoring Application
#### Core Features:
- Multi-City Weather Tracking
  - Real-time monitoring of three cities
  - Data collection every 5 minutes
  - Tracks:
    - Current temperature
    - "Feels like" temperature
    - Dominant weather condition

#### Data Analysis:
- Daily Data Aggregation and Rollups:
  - Maximum temperature
  - Minimum temperature
  - Average temperature
  - Per-city analysis

#### Advanced Features:
- Customizable Alert System
  - User-defined temperature thresholds for each city
  - Automated notification system
  - City-specific alert configurations

#### Performance Optimization:
- Query optimization using:
  - Database Views
  - TimescaleDB time_bucket function
  - Efficient data rollup strategies

## Getting Started

There are two ways to run these applications:

### Option 1: Using Docker Compose

1. Clone the repository:
```bash
git clone https://github.com/TejasThombare20/RuleEngine-WeatherMonitor.git
cd zeotap
```

2. Start all services using Docker Compose:
```bash
docker-compose up -d
```

3. Access the applications:
   - Rule Engine: http://localhost:3001
   - Weather Monitor: http://localhost:3002

### Option 2: Using Pre-built Docker Image

1. Pull the image:
```bash
docker pull tejasthombare/zeotap
```

2. Run the container:
```bash
docker run -d \
  -p 3001:3000 \
  -p 3002:3000 \
  -p 8000:8000 \
  -p 9000:9000 \
  tejasthombare/zeotap
```

3. Access the applications:
   - Rule Engine Frontend: http://localhost:3001
   - Rule Engine Backend: http://localhost:8000
   - Weather Monitor Frontend: http://localhost:3002
   - Weather Monitor Backend: http://localhost:9000

## Application Ports

### Rule Engine
- Frontend: 3001
- Backend API: 8000

### Weather Monitor
- Frontend: 3002
- Backend API: 9000

## Project Structure

```
.
├── rule_engine/
│   ├── client/          # Next.js frontend
│   └── server/          # Golang backend
├── weather_monitoring/
│   ├── client/          # Next.js frontend
│   └── server/          # Golang backend
├── Dockerfile
├── docker-compose.yml
└── README.md
```



