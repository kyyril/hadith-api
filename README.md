# Hadith API

API Hadis dengan terjemahan dari 9 perawi, data dalam file JSON lokal.

## Features

- Get hadiths by narrator (`/api/v1/hadis/:slug`)
- Get a specific hadith by narrator and number (`/api/v1/hadis/:slug/:number`)
- Get list of available narrators (`/api/v1/narrators`)
- Swagger documentation
- Pagination and search support

## API Endpoints

### Get Available Narrators

```
GET /api/v1/narrators
```

Returns a list of all available hadith narrators.

### Get Hadiths by Narrator

```
GET /api/v1/hadis/:slug
```

Returns all hadiths from a specific narrator.

Query parameters:
- `page`: Page number for pagination (default: 1)
- `limit`: Number of hadiths per page (default: 10, max: 100)
- `q`: Search query to filter hadiths

### Get Hadith by Number

```
GET /api/v1/hadis/:slug/:number
```

Returns a specific hadith from a narrator by its number.

## Data Format

Each hadith is stored in the following format:

```json
{
  "number": 1,
  "arab": "حَدَّثَنَا أَبُو بَكْرِ بْنُ أَبِي شَيْبَةَ...",
  "id": "Telah menceritakan kepada kami Abu Bakar bin Abu Syaibah..."
}
```

## Local Development

```bash
# Clone repository
git clone https://github.com/kyyril/hadiths-restapi.git

# Navigate to the project
cd hadiths-restapi

# Run the server
go run main.go

# The server will start on port 8080 by default
```

## Deployment

This API is designed to be deployable on Vercel.

## Documentation

Swagger documentation is available at `/swagger/index.html` when running the server.