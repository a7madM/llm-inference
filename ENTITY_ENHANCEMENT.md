# Entity Enhancement Service

## Overview

The Entity Enhancement Service processes arrays of named entities extracted from other services to improve their quality and remove duplicates using Large Language Models (LLMs). This service is particularly useful when working with noisy or low-quality entity extraction results from other NLP tools.

## Features

- **Quality Improvement**: Corrects typos, standardizes formatting, and improves entity quality
- **Deduplication**: Removes exact duplicates and near-duplicates (e.g., "John Smith" vs "john smith")
- **Entity Type Awareness**: Tailors enhancement based on entity type (person, location, organization, etc.)
- **Intelligent Filtering**: Removes invalid or low-quality entities that don't match the expected type
- **LLM-Powered**: Uses advanced language models for semantic understanding and enhancement

## API Endpoint

### POST `/api/v1/enhance-entities`

Enhances an array of named entities by improving quality and removing duplicates.

#### Request Body

```json
{
  "entities": [
    "john smith",
    "JOHN SMITH",
    "Jon Smith",
    "New York",
    "new york city", 
    "NYC",
    "Apple Inc",
    "Apple Inc.",
    "invalid_entity_123"
  ],
  "entity_type": "person"
}
```

#### Response

```json
{
  "original_entities": [
    "john smith",
    "JOHN SMITH", 
    "Jon Smith",
    "New York",
    "new york city",
    "NYC",
    "Apple Inc",
    "Apple Inc.",
    "invalid_entity_123"
  ],
  "enhanced_entities": [
    "John Smith",
    "New York City",
    "Apple Inc"
  ],
  "entity_type": "person",
  "processed_count": 9,
  "removed_count": 6,
  "thinking": "Analysis of entity enhancement process..."
}
```

## Usage Examples

### Example 1: Person Names

**Request:**
```bash
curl -X POST http://localhost:8090/api/v1/enhance-entities \
  -H "Content-Type: application/json" \
  -d '{
    "entities": [
      "john doe",
      "JOHN DOE",
      "Jon Doe", 
      "Jane Smith",
      "j. smith",
      "Dr. Jane Smith",
      "jane smith md"
    ],
    "entity_type": "person"
  }'
```

**Response:**
```json
{
  "original_entities": ["john doe", "JOHN DOE", "Jon Doe", "Jane Smith", "j. smith", "Dr. Jane Smith", "jane smith md"],
  "enhanced_entities": ["John Doe", "Dr. Jane Smith"],
  "entity_type": "person",
  "processed_count": 7,
  "removed_count": 5
}
```

### Example 2: Organizations

**Request:**
```bash
curl -X POST http://localhost:8090/api/v1/enhance-entities \
  -H "Content-Type: application/json" \
  -d '{
    "entities": [
      "microsoft corp",
      "Microsoft Corporation", 
      "MSFT",
      "Google LLC",
      "google inc",
      "Alphabet Inc",
      "not_a_company_123"
    ],
    "entity_type": "organization"
  }'
```

**Response:**
```json
{
  "original_entities": ["microsoft corp", "Microsoft Corporation", "MSFT", "Google LLC", "google inc", "Alphabet Inc", "not_a_company_123"],
  "enhanced_entities": ["Microsoft Corporation", "Google LLC", "Alphabet Inc"],
  "entity_type": "organization", 
  "processed_count": 7,
  "removed_count": 4
}
```

### Example 3: Locations

**Request:**
```bash
curl -X POST http://localhost:8090/api/v1/enhance-entities \
  -H "Content-Type: application/json" \
  -d '{
    "entities": [
      "new york",
      "New York City",
      "NYC", 
      "manhattan",
      "los angeles",
      "LA",
      "california",
      "random_text_456"
    ],
    "entity_type": "location"
  }'
```

**Response:**
```json
{
  "original_entities": ["new york", "New York City", "NYC", "manhattan", "los angeles", "LA", "california", "random_text_456"],
  "enhanced_entities": ["New York City", "Manhattan", "Los Angeles", "California"],
  "entity_type": "location",
  "processed_count": 8,
  "removed_count": 4
}
```

## Enhancement Logic

The service uses sophisticated LLM prompts to:

1. **Standardize Formatting**: Convert to proper case and consistent formatting
2. **Remove Duplicates**: Identify and merge similar entities (semantic similarity)
3. **Quality Filtering**: Remove entities that don't match the expected type
4. **Error Correction**: Fix common typos and formatting issues
5. **Consolidation**: Merge different representations of the same entity

## Supported Entity Types

- `person` - Individual names
- `organization` - Company and organization names  
- `location` - Geographic locations (cities, countries, landmarks)
- `event` - Event names and titles
- `product` - Product names and models
- `misc` - Miscellaneous entities

## Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `original_entities` | `[]string` | The input array of entities |
| `enhanced_entities` | `[]string` | The processed and enhanced entities |
| `entity_type` | `string` | The type of entities being processed |
| `processed_count` | `int` | Number of original entities processed |
| `removed_count` | `int` | Number of entities removed during enhancement |
| `thinking` | `string` | LLM reasoning process (optional) |

## Error Handling

### Validation Errors (400 Bad Request)

```json
{
  "error": "Key: 'EntityEnhancementRequest.Entities' Error:Field validation for 'Entities' failed on the 'required' tag"
}
```

### Processing Errors (500 Internal Server Error)

```json
{
  "error": "Failed to enhance entities"
}
```

## Configuration

The service uses the same configuration as other LLM endpoints:

```env
OLLAMA_URL=http://localhost:11434
MODEL_NAME=deepseek-r1:1.5b
PORT=8090
```

## Performance Considerations

- **Batch Size**: Optimal performance with 10-50 entities per request
- **Entity Type**: More specific types (e.g., "person") perform better than generic types
- **Quality**: Higher quality input entities result in better enhancement results
- **Caching**: Consider caching results for frequently processed entity sets

## Integration Examples

### Python Integration

```python
import requests

def enhance_entities(entities, entity_type, base_url="http://localhost:8090"):
    response = requests.post(
        f"{base_url}/api/v1/enhance-entities",
        json={
            "entities": entities,
            "entity_type": entity_type
        }
    )
    return response.json()

# Usage
entities = ["john doe", "JOHN DOE", "Jane Smith"]
result = enhance_entities(entities, "person")
print(f"Enhanced: {result['enhanced_entities']}")
```

### JavaScript Integration

```javascript
async function enhanceEntities(entities, entityType, baseUrl = "http://localhost:8090") {
    const response = await fetch(`${baseUrl}/api/v1/enhance-entities`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            entities: entities,
            entity_type: entityType
        })
    });
    return await response.json();
}

// Usage
const entities = ["apple inc", "Apple Inc.", "AAPL"];
const result = await enhanceEntities(entities, "organization");
console.log("Enhanced:", result.enhanced_entities);
```

### Go Integration

```go
type EntityRequest struct {
    Entities   []string `json:"entities"`
    EntityType string   `json:"entity_type"`
}

func enhanceEntities(entities []string, entityType string) (*EntityEnhancementResponse, error) {
    reqBody := EntityRequest{
        Entities:   entities,
        EntityType: entityType,
    }
    
    jsonData, _ := json.Marshal(reqBody)
    resp, err := http.Post("http://localhost:8090/api/v1/enhance-entities", 
                          "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result EntityEnhancementResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}
```

## Best Practices

1. **Batch Processing**: Group entities of the same type for better efficiency
2. **Type Specificity**: Use specific entity types rather than generic ones
3. **Pre-filtering**: Remove obviously invalid entities before enhancement
4. **Result Validation**: Verify enhanced results match your quality requirements
5. **Caching**: Cache results for frequently processed entity sets
6. **Error Handling**: Implement proper error handling for failed requests

## Troubleshooting

### Common Issues

1. **Empty Results**: Check if entity_type matches the actual entities
2. **Poor Quality**: Ensure input entities are reasonably formatted
3. **Timeout**: Large batches may timeout; reduce batch size
4. **Model Issues**: Verify Ollama is running and model is available

### Debug Information

Enable debug logging to see LLM reasoning:

```bash
curl -X POST http://localhost:8090/api/v1/enhance-entities \
  -H "Content-Type: application/json" \
  -d '{"entities": ["test"], "entity_type": "person"}' | jq '.thinking'
```