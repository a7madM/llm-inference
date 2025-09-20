import os
from fastapi import FastAPI
from pydantic import BaseModel
from typing import List
import requests
import json
from time import time

# Ollama API settings
OLLAMA_URL = os.getenv("OLLAMA_URL", "http://localhost:11434")
MODEL_NAME = os.getenv("MODEL_NAME", "deepseek-r1:1.5b")
API_URL = f"{OLLAMA_URL}/api/generate"

app = FastAPI(
    title="Local Multilingual NER Service",
    description="Extracts Persons, Locations, Organizations, Events from text (Arabic, English, German) using Mistral-7B. Preserves original language.",
    version="1.2.0",
)

# ---- Input / Output Schemas ----
class InputText(BaseModel):
    text: str

class Entities(BaseModel):
    persons: List[str] = []
    locations: List[str] = []
    organizations: List[str] = []
    events: List[str] = []

# ---- Helper: Safe JSON Parsing ----
def safe_json_parse(raw_output: str) -> dict:
    """
    Safely extract JSON from model output.
    Fallbacks to empty schema if parsing fails.
    """
    try:
        return json.loads(raw_output)
    except json.JSONDecodeError:
        start = raw_output.find("{")
        end = raw_output.rfind("}") + 1
        if start != -1 and end != -1:
            try:
                return json.loads(raw_output[start:end])
            except Exception:
                pass
    return {"persons": [], "locations": [], "organizations": [], "events": []}

# ---- API Endpoint ----
@app.post("/api/v1/entities", response_model=Entities)
def extract_entities(input_data: InputText):
    print(f"Received text for NER: {input_data.text[:50]}...")
    timeNow = time()
    prompt = f"""
You are an information extraction system. 
Extract named entities from the following text **without translating it**.
The text may be in Arabic, English, or German.

Rules:
1. Person names must consist of at least **two words**.
2. Events must be **phrases**, not single words.
3. Return only valid JSON with this structure:
{{
  "persons": ["string"],
  "locations": ["string"],
  "organizations": ["string"],
  "events": ["string"]
}}
4. Do not translate the text. Do not include explanations, extra text, or markdown. 
5. Keep all names and words in the original language.

Text: {input_data.text}
"""

    # Call Ollama API
    response = requests.post(
        API_URL,
        json={"model": MODEL_NAME, "prompt": prompt, "stream": False},
    )

    if response.status_code != 200:
        return { "response": "Error: Unable to reach Ollama API." }

    output_text = response.json().get("response", "").strip()
    data = safe_json_parse(output_text)

    # ---- Post-processing ----
    # Keep only persons with at least 2 words
    data["persons"] = [p for p in data.get("persons", []) if len(p.strip().split()) >= 2]

    # Keep only events with more than 1 word
    data["events"] = [t for t in data.get("events", []) if len(t.strip().split()) > 1]

    timeElapsed = time() - timeNow
    print(f"Extraction completed in {timeElapsed:.2f} seconds.")
    return Entities(**data)


# ---- Sentiment Schema ----
class Sentiment(BaseModel):
    sentiment: str  # "positive", "negative", "neutral"
    confidence: float  # 0.0 to 1.0
@app.post("/api/v1/sentiment", response_model=Sentiment)
def analyze_sentiment(input_data: InputText):
        print(f"Analyzing sentiment for: {input_data.text[:50]}...")
        timeNow = time()
        
        prompt = f"""
    Analyze the sentiment of the following text. The text may be in Arabic, English, or German.
    Return only valid JSON with this structure:
    {{
      "sentiment": "positive|negative|neutral",
      "confidence": 0.0-1.0
    }}

    Text: {input_data.text}
    """

        response = requests.post(
            API_URL,
            json={"model": MODEL_NAME, "prompt": prompt, "stream": False},
        )

        if response.status_code != 200:
            return {"sentiment": "neutral", "confidence": 0.0}

        try:
            output = response.json().get("response", "").strip()
            data = json.loads(output)
            timeElapsed = time() - timeNow
            print(f"Sentiment analysis completed in {timeElapsed:.2f} seconds.")
            return Sentiment(**data)
        except:
            return {"sentiment": "neutral", "confidence": 0.0}