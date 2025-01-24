package openaiassistantsgo

import (
  "os"
  "testing"
  "time"
  "fmt"

  "github.com/stretchr/testify/assert"
  "github.com/joho/godotenv"
)

func TestAssistantStructuredOutputs(t *testing.T) {
  // Load environment variables

  // Load .env.local file if it exists
  if err := godotenv.Load(".env.local"); err != nil {
    t.Fatalf("No .env.local file found: %v", err)
  }
  apiKey := os.Getenv("OPENAI_API_KEY")
  if apiKey == "" {
    t.Fatal("OPENAI_API_KEY environment variable not set")
  }

  // Define schema file path
  //schemaFilePath := "./structured_output_schema.json"

  // Step 1: Create a Vector Store
  vectorStoreID, err := CreateVectorStore("Test Vector Store")
  assert.NoError(t, err, "Error creating vector store")
  assert.NotEmpty(t, vectorStoreID, "Vector store ID should not be empty")
  t.Logf("Vector Store created: %s", vectorStoreID)

  // Step 2: Create an Assistant
  assistantID, err := CreateAssistant("You create value proposals.", vectorStoreID, "")
  assert.NoError(t, err, "Error creating assistant")
  assert.NotEmpty(t, assistantID, "Assistant ID should not be empty")
  t.Logf("Assistant created: %s", assistantID)

  // Step 3: Create a Thread
  threadID, err := CreateThread()
  assert.NoError(t, err, "Error creating thread")
  assert.NotEmpty(t, threadID, "Thread ID should not be empty")
  t.Logf("Thread created: %s", threadID)

  // Step 4: Send a Message to the Thread
  messageID, err := CreateMessage(threadID, "Generate a structured json output with result 'success' and key-value data.", "user")
  assert.NoError(t, err, "Error creating message")
  assert.NotEmpty(t, messageID, "Message ID should not be empty")
  t.Logf("Message sent: %s", messageID)

  // Step 5: Run the Thread with the Assistant
  runID, err := CreateRun(assistantID, threadID)
  assert.NoError(t, err, "Error creating Run.")
  assert.NotNil(t, runID, "runID should not be nil.")
  t.Logf("Run created: %s", runID)

  // Step 6: Retrieve the Run Status
  for {
    time.Sleep(time.Second)

    runStatus, err := GetRunStatus(runID, threadID)
    if err != nil {
      t.Fatal(err)
    }
    assert.NoError(t, err, "Error retrieving run status.")
    assert.NotNil(t, runStatus, "runStatus should not be nil.")

    if runStatus == "completed" {
      t.Logf("Run status: %s", runStatus)
      break
    }

    t.Logf("Run status: %s", runStatus)

  }

  // Step 7: Get the Structured Response
  messages, err := ListMessages(threadID)
  if err != nil {
    t.Fatalf("%v...%v", "error", err)
    return
  }
  assert.NoError(t, err, "Error retrieving message")
  assert.NotNil(t, messages, "Messages should not be nil")

  // Grab first message
  assistantRespID := messages.FirstID
  message, err := GetMessage(threadID, assistantRespID)
  if err != nil {
    t.Fatalf("%v...%v", "error", err)
    return
  }
  assert.NoError(t, err, "Error retrieving message")
  assert.NotNil(t, message, "Message should not be nil")

  t.Logf("Message Content%s", message.Content[0].Text.Value)
  fmt.Printf("Message Content%s", message.Content[0].Text.Value)

  // Step 8: Cleanup Resources
  err = DeleteThread(threadID)
  assert.NoError(t, err, "Error deleting thread")
  t.Logf("Thread deleted: %s", threadID)

  err = DeleteAssistant(assistantID)
  assert.NoError(t, err, "Error deleting assistant.")
  t.Logf("Vector store deleted: %s", vectorStoreID)

  err = DeleteVectorStore(vectorStoreID)
  assert.NoError(t, err, "Error deleting vector store")
  t.Logf("Vector store deleted: %s", vectorStoreID)
}
