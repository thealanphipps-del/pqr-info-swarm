const express = require('express');
const cors = require('cors');
const fs = require('fs');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 3001;

app.use(cors());
app.use(express.json());

// Function to parse the mock patent text
function parsePatentText(text) {
  const result = {
    title: '',
    patentNumber: '',
    inventor: '',
    abstract: '',
    description: '',
    claims: []
  };

  const lines = text.split('\n');
  let currentSection = '';

  for (const line of lines) {
    if (line.startsWith('Title: ')) {
      result.title = line.replace('Title: ', '').trim();
    } else if (line.startsWith('Patent Number: ')) {
      result.patentNumber = line.replace('Patent Number: ', '').trim();
    } else if (line.startsWith('Inventor: ')) {
      result.inventor = line.replace('Inventor: ', '').trim();
    } else if (line.startsWith('Abstract: ')) {
      result.abstract = line.replace('Abstract: ', '').trim();
      currentSection = 'abstract';
    } else if (line.startsWith('Description:')) {
      currentSection = 'description';
    } else if (line.startsWith('Claims:')) {
      currentSection = 'claims';
    } else {
      if (currentSection === 'description') {
        if (line.trim()) {
          result.description += line + '\n';
        }
      } else if (currentSection === 'claims') {
        if (line.trim()) {
          result.claims.push(line.trim());
        }
      } else if (currentSection === 'abstract') {
        if (line.trim()) {
          result.abstract += ' ' + line.trim();
        }
      }
    }
  }

  result.description = result.description.trim();
  return result;
}

// Endpoint to fetch and parse the patent context
app.get('/api/context', (req, res) => {
  try {
    const patentPath = path.join(__dirname, 'mock_patent.txt');
    const rawText = fs.readFileSync(patentPath, 'utf8');
    const parsedContext = parsePatentText(rawText);
    
    res.json({
      success: true,
      data: parsedContext
    });
  } catch (error) {
    console.error('Error reading or parsing patent file:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to fetch context'
    });
  }
});

// Also provide a raw text endpoint just in case
app.get('/api/context/raw', (req, res) => {
  try {
    const patentPath = path.join(__dirname, 'mock_patent.txt');
    const rawText = fs.readFileSync(patentPath, 'utf8');
    
    res.json({
      success: true,
      data: rawText
    });
  } catch (error) {
    console.error('Error reading patent file:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to fetch raw context'
    });
  }
});

app.listen(PORT, () => {
  console.log(`Context fetching service running on port ${PORT}`);
});
