# PQR Ticketing System - Windows Agent Memory Test
# Usage: .\test-agent-memory.ps1 -BaseUrl http://localhost:8080 -AgentId test-agent-001

param(
    [string]$BaseUrl = "http://localhost:8080",
    [string]$AgentId = "test-agent-001"
)

Write-Host "PQR Ticketing System - Agent Memory Test" -ForegroundColor Green
Write-Host "==========================================" 
Write-Host "Base URL: $BaseUrl"
Write-Host "Agent ID: $AgentId`n"

# Helper function
function Invoke-TicketAPI {
    param(
        [string]$Method,
        [string]$Endpoint,
        [hashtable]$Body
    )
    
    $uri = "$BaseUrl/REST/2.0$Endpoint"
    $headers = @{ "Content-Type" = "application/json" }
    
    if ($Body) {
        $bodyJson = $Body | ConvertTo-Json -Depth 10
        return Invoke-RestMethod -Uri $uri -Method $Method -Headers $headers -Body $bodyJson -ErrorAction Stop
    } else {
        return Invoke-RestMethod -Uri $uri -Method $Method -Headers $headers -ErrorAction Stop
    }
}

# 1. Health Check
Write-Host "1. Health Check..." -ForegroundColor Cyan
$health = Invoke-TicketAPI -Method GET -Endpoint "/health"
$health | ConvertTo-Json | Write-Host
Write-Host

# 2. Initialize Schema
Write-Host "2. Initialize Schema..." -ForegroundColor Cyan
$init = Invoke-TicketAPI -Method POST -Endpoint "/init"
$init | ConvertTo-Json | Write-Host
Write-Host

# 3. Create Ticket
Write-Host "3. Creating Ticket..." -ForegroundColor Cyan
$ticketBody = @{
    Subject = "Agent Working Memory"
    Queue = "processing"
    Text = "Initial task content"
    AgentID = $AgentId
    Layer = 2
    Intent = @{
        task = "test"
        priority = "high"
    }
}
$ticket = Invoke-TicketAPI -Method POST -Endpoint "/ticket" -Body $ticketBody
$ticketId = $ticket.id
Write-Host "Ticket ID: $ticketId" -ForegroundColor Yellow
Write-Host

# 4. Store Memory
Write-Host "4. Storing Agent Memory..." -ForegroundColor Cyan
$memoryBody = @{
    memory_type = "context"
    data = @{
        status = "processing"
        items_processed = 5
        items_total = 10
        current_item = "data_point_5"
    }
    relevance_score = 0.95
}
$memory = Invoke-TicketAPI -Method POST -Endpoint "/agent/$AgentId/memory/$ticketId" -Body $memoryBody
$memory | ConvertTo-Json | Write-Host
Write-Host

# 5. Retrieve Memory
Write-Host "5. Retrieving Agent Memory..." -ForegroundColor Cyan
$retrieved = Invoke-TicketAPI -Method GET -Endpoint "/agent/$AgentId/memory/$ticketId?type=context"
$retrieved | ConvertTo-Json | Write-Host
Write-Host

# 6. Get Ticket Details
Write-Host "6. Getting Ticket Details..." -ForegroundColor Cyan
$ticketDetails = Invoke-TicketAPI -Method GET -Endpoint "/ticket/$ticketId"
$ticketDetails | ConvertTo-Json -Depth 10 | Write-Host
Write-Host

# 7. Store Knowledge Memory
Write-Host "7. Storing Knowledge Memory..." -ForegroundColor Cyan
$knowledgeBody = @{
    memory_type = "knowledge"
    data = @{
        patterns = @("pattern_a", "pattern_b")
        confidence = 0.87
    }
    relevance_score = 0.85
}
$knowledge = Invoke-TicketAPI -Method POST -Endpoint "/agent/$AgentId/memory/$ticketId" -Body $knowledgeBody
$knowledge | ConvertTo-Json | Write-Host
Write-Host

# 8. Update Ticket Status
Write-Host "8. Updating Ticket Status..." -ForegroundColor Cyan
$updateBody = @{
    Status = "PROCESSING"
    Title = "Updated: Memory Storage Test"
}
$updated = Invoke-TicketAPI -Method PUT -Endpoint "/ticket/$ticketId" -Body $updateBody
$updated | ConvertTo-Json | Write-Host
Write-Host

# 9. Get Audit Trail
Write-Host "9. Getting Audit Trail..." -ForegroundColor Cyan
$audit = Invoke-TicketAPI -Method GET -Endpoint "/ticket/$ticketId/audit"
$audit | ConvertTo-Json -Depth 10 | Write-Host
Write-Host

# 10. Get Agent Context
Write-Host "10. Getting Agent Context..." -ForegroundColor Cyan
$context = Invoke-TicketAPI -Method GET -Endpoint "/agent/$AgentId/context"
$context | ConvertTo-Json -Depth 10 | Write-Host
Write-Host

Write-Host "==========================================" -ForegroundColor Green
Write-Host "Test Complete!" -ForegroundColor Green
Write-Host "Ticket ID: $ticketId" -ForegroundColor Yellow
Write-Host "Use this to test further operations:" -ForegroundColor Cyan
Write-Host "  - GET /REST/2.0/ticket/$ticketId"
Write-Host "  - GET /REST/2.0/agent/$AgentId/memory/$ticketId"
Write-Host "  - PUT /REST/2.0/ticket/$ticketId"


