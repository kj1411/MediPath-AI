# Design Document: MediPath AI

## Overview

MediPath AI is a patient-facing AI care navigation assistant designed to help patients understand medical information provided after doctor consultations. The system operates under strict safety constraints to ensure it never provides diagnostic advice, prescribes medications, or contradicts medical professionals.

### Design Principles

1. **Safety First**: All responses must pass through guardrails that prevent harmful medical advice
2. **Doctor Authority**: The system always defers to and supports doctor decisions
3. **Privacy by Design**: No persistent storage of Protected Health Information (PHI)
4. **Graceful Degradation**: System operates in minimal-input mode when documents are unavailable
5. **Transparency**: Clear disclaimers and confidence indicators in all responses
6. **Accessibility**: Support for multiple languages, voice input, and accessibility standards

### System Boundaries

**In Scope:**
- Explaining existing diagnoses in simple language
- Clarifying prescribed medication usage
- Interpreting lab reports with context
- Providing follow-up care guidance
- Offering general preventive health education
- Detecting red-flag symptoms for triage

**Out of Scope:**
- Diagnosing diseases or conditions
- Prescribing medications or treatments
- Recommending treatment alternatives
- Contradicting doctor decisions
- Emergency medical response
- Long-term PHI storage

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Web UI     │  │  Mobile UI   │  │ Voice Input  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      API Gateway Layer                       │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Request Router │ Auth │ Rate Limiter │ Encryption  │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Processing Pipeline                       │
│                                                               │
│  ┌──────────────┐      ┌──────────────┐                     │
│  │ OCR Engine   │──────▶│ NER Engine   │                     │
│  └──────────────┘      └──────────────┘                     │
│         │                      │                              │
│         └──────────┬───────────┘                              │
│                    ▼                                          │
│         ┌──────────────────────┐                             │
│         │  Entity Validator    │                             │
│         └──────────────────────┘                             │
│                    │                                          │
│                    ▼                                          │
│         ┌──────────────────────┐                             │
│         │ Context Builder      │                             │
│         └──────────────────────┘                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                        AI Layer                              │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Explanation Engine (LLM)                 │   │
│  │  ┌────────────────────────────────────────────────┐  │   │
│  │  │  Prompt Constructor │ Response Generator       │  │   │
│  │  └────────────────────────────────────────────────┘  │   │
│  └──────────────────────────────────────────────────────┘   │
│                    │                                          │
│                    ▼                                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Guardrail System                         │   │
│  │  ┌────────────────────────────────────────────────┐  │   │
│  │  │ Safety Rules │ Confidence Check │ Abstention   │  │   │
│  │  └────────────────────────────────────────────────┘  │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Knowledge Layer                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Medical    │  │  Red Flag    │  │  Medication  │      │
│  │  Guidelines  │  │   Symptoms   │  │   Database   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

### Architectural Patterns

1. **Pipeline Architecture**: Document processing follows a sequential pipeline (OCR → NER → Validation → Context Building → Explanation)
2. **Layered Architecture**: Clear separation between presentation, business logic, AI processing, and data layers
3. **Guardrail Pattern**: All AI-generated responses pass through safety validation before delivery
4. **Stateless Sessions**: No persistent storage of PHI; all processing happens in ephemeral sessions
5. **Fail-Safe Design**: System defaults to safe responses (abstention + doctor consultation) when uncertain

## Components and Interfaces

### Frontend Components

#### Web/Mobile UI
**Responsibilities:**
- Render chat interface for patient interactions
- Handle document uploads (PDF, images, text)
- Display formatted responses with disclaimers
- Provide language selection
- Show processing status and feedback

**Interfaces:**
```typescript
interface UIComponent {
  renderChatMessage(message: Message): void
  handleDocumentUpload(file: File): Promise<UploadResult>
  displayResponse(response: ExplanationResponse): void
  showError(error: ErrorMessage): void
  setLanguage(language: LanguageCode): void
}

interface Message {
  id: string
  sender: 'patient' | 'system'
  content: string
  timestamp: Date
  disclaimer?: string
}

interface UploadResult {
  success: boolean
  documentId?: string
  error?: string
}
```

#### Voice Input Module (Optional)
**Responsibilities:**
- Capture voice input from patient
- Convert speech to text
- Handle voice input errors

**Interfaces:**
```typescript
interface VoiceInput {
  startRecording(): void
  stopRecording(): Promise<string>
  convertSpeechToText(audio: AudioData): Promise<string>
}
```

### API Gateway Layer

#### Request Router
**Responsibilities:**
- Route requests to appropriate processing pipelines
- Handle authentication and session management
- Apply rate limiting
- Encrypt data in transit

**Interfaces:**
```typescript
interface APIGateway {
  routeRequest(request: PatientRequest): Promise<Response>
  authenticateSession(sessionToken: string): Promise<SessionInfo>
  applyRateLimit(sessionId: string): boolean
  encryptData(data: any): EncryptedData
  decryptData(encrypted: EncryptedData): any
}

interface PatientRequest {
  sessionId: string
  type: 'document_upload' | 'question' | 'follow_up'
  payload: DocumentPayload | QuestionPayload
  language: LanguageCode
}
```

### Processing Pipeline

#### OCR Engine
**Responsibilities:**
- Extract text from PDF documents
- Extract text from image documents (JPEG, PNG)
- Handle various document qualities
- Detect document language

**Interfaces:**
```typescript
interface OCREngine {
  extractTextFromPDF(pdf: PDFDocument): Promise<ExtractedText>
  extractTextFromImage(image: ImageDocument): Promise<ExtractedText>
  detectLanguage(text: string): LanguageCode
}

interface ExtractedText {
  content: string
  confidence: number
  language: LanguageCode
  metadata: DocumentMetadata
}
```

#### NER Engine (Named Entity Recognition)
**Responsibilities:**
- Identify medical entities in extracted text
- Extract diagnoses, medications, lab values, procedures
- Tag entity types and relationships
- Handle medical terminology variations

**Interfaces:**
```typescript
interface NEREngine {
  extractEntities(text: string): Promise<MedicalEntities>
  validateEntityTypes(entities: MedicalEntities): ValidationResult
}

interface MedicalEntities {
  diagnoses: Diagnosis[]
  medications: Medication[]
  labResults: LabResult[]
  procedures: Procedure[]
  followUpInstructions: Instruction[]
}

interface Diagnosis {
  name: string
  icdCode?: string
  confidence: number
  context: string
}

interface Medication {
  name: string
  dosage?: string
  frequency?: string
  duration?: string
  purpose?: string
}

interface LabResult {
  testName: string
  value: string
  unit?: string
  normalRange?: string
  flag?: 'high' | 'low' | 'normal'
}
```

#### Entity Validator
**Responsibilities:**
- Validate extracted entities against medical knowledge base
- Check entity completeness and consistency
- Flag ambiguous or uncertain entities

**Interfaces:**
```typescript
interface EntityValidator {
  validateEntities(entities: MedicalEntities): ValidationResult
  checkCompleteness(entities: MedicalEntities): CompletenessReport
  resolveAmbiguities(entities: MedicalEntities): ResolvedEntities
}

interface ValidationResult {
  isValid: boolean
  validatedEntities: MedicalEntities
  warnings: ValidationWarning[]
  errors: ValidationError[]
}
```

#### Context Builder
**Responsibilities:**
- Combine extracted entities with patient questions
- Build context for LLM prompt
- Maintain conversation history within session
- Apply minimal-input mode logic when no documents present

**Interfaces:**
```typescript
interface ContextBuilder {
  buildContext(entities: MedicalEntities, question: string, history: ConversationHistory): PromptContext
  buildMinimalInputContext(question: string): PromptContext
  addToHistory(exchange: ConversationExchange): void
  clearHistory(sessionId: string): void
}

interface PromptContext {
  mode: 'document_based' | 'minimal_input'
  extractedEntities?: MedicalEntities
  patientQuestion: string
  conversationHistory: ConversationExchange[]
  language: LanguageCode
}

interface ConversationExchange {
  patientMessage: string
  systemResponse: string
  timestamp: Date
}
```

### AI Layer

#### Explanation Engine (LLM)
**Responsibilities:**
- Generate patient-friendly explanations
- Ground responses in extracted entities
- Maintain conversational tone
- Apply language-specific generation
- Include appropriate disclaimers

**Interfaces:**
```typescript
interface ExplanationEngine {
  generateExplanation(context: PromptContext): Promise<RawExplanation>
  constructPrompt(context: PromptContext): string
  applyDisclaimers(explanation: string): string
}

interface RawExplanation {
  content: string
  confidence: number
  entitiesReferenced: string[]
  requiresGuardrailCheck: boolean
}
```

#### Guardrail System
**Responsibilities:**
- Validate responses against safety rules
- Detect prohibited content (diagnoses, prescriptions, contradictions)
- Check response confidence levels
- Trigger abstention when appropriate
- Generate safe alternative responses

**Interfaces:**
```typescript
interface GuardrailSystem {
  validateResponse(response: RawExplanation, context: PromptContext): GuardrailResult
  checkSafetyRules(response: string): SafetyCheckResult
  checkConfidence(confidence: number): boolean
  generateAbstentionResponse(reason: AbstentionReason): string
}

interface GuardrailResult {
  approved: boolean
  finalResponse: string
  violations: SafetyViolation[]
  action: 'approve' | 'modify' | 'abstain'
}

interface SafetyViolation {
  rule: SafetyRule
  severity: 'critical' | 'warning'
  detectedContent: string
}

enum SafetyRule {
  NO_DIAGNOSIS = 'no_diagnosis',
  NO_PRESCRIPTION = 'no_prescription',
  NO_TREATMENT_ALTERNATIVES = 'no_treatment_alternatives',
  NO_CONTRADICTION = 'no_contradiction',
  CONFIDENCE_THRESHOLD = 'confidence_threshold'
}

enum AbstentionReason {
  LOW_CONFIDENCE = 'low_confidence',
  SAFETY_VIOLATION = 'safety_violation',
  INSUFFICIENT_CONTEXT = 'insufficient_context',
  AMBIGUOUS_ENTITIES = 'ambiguous_entities'
}
```

#### Red Flag Detector
**Responsibilities:**
- Check patient-described symptoms against red flag patterns
- Determine urgency level
- Generate appropriate triage recommendations

**Interfaces:**
```typescript
interface RedFlagDetector {
  checkSymptoms(symptoms: string): RedFlagResult
  determineUrgency(symptoms: string): UrgencyLevel
}

interface RedFlagResult {
  hasRedFlags: boolean
  detectedSymptoms: RedFlagSymptom[]
  urgencyLevel: UrgencyLevel
  recommendation: string
}

interface RedFlagSymptom {
  symptom: string
  category: 'emergency' | 'urgent' | 'routine'
  description: string
}

enum UrgencyLevel {
  EMERGENCY = 'emergency',        // Call 911 / Emergency services
  URGENT = 'urgent',              // Contact doctor within 24 hours
  ROUTINE = 'routine',            // Schedule regular appointment
  INFORMATIONAL = 'informational' // General education, no immediate action
}
```

### Knowledge Layer

#### Medical Guidelines Database
**Responsibilities:**
- Store public medical guidelines and educational content
- Provide medication information
- Supply preventive health recommendations

**Interfaces:**
```typescript
interface MedicalKnowledgeBase {
  getMedicationInfo(medicationName: string): MedicationInfo
  getPreventiveAdvice(condition: string): PreventiveAdvice[]
  getGeneralEducation(topic: string): EducationalContent
}

interface MedicationInfo {
  name: string
  commonUses: string[]
  sideEffects: string[]
  generalInstructions: string
  source: string
}

interface PreventiveAdvice {
  category: string
  recommendation: string
  evidence: string
  source: string
}
```

#### Red Flag Symptom Database
**Responsibilities:**
- Store rule-based red flag symptom patterns
- Provide urgency classifications
- Supply triage recommendations

**Interfaces:**
```typescript
interface RedFlagDatabase {
  getRedFlagPatterns(): RedFlagPattern[]
  getTriageGuidelines(urgency: UrgencyLevel): TriageGuideline
}

interface RedFlagPattern {
  symptomKeywords: string[]
  urgencyLevel: UrgencyLevel
  recommendation: string
}
```

### Security Layer

#### Encryption Service
**Responsibilities:**
- Encrypt data in transit (TLS)
- Encrypt data at rest during processing
- Secure key management

**Interfaces:**
```typescript
interface EncryptionService {
  encryptInTransit(data: any): EncryptedData
  encryptAtRest(data: any): EncryptedData
  decrypt(encrypted: EncryptedData): any
  rotateKeys(): void
}
```

#### Session Manager
**Responsibilities:**
- Create isolated processing sessions
- Delete session data after completion
- Ensure no PHI persistence

**Interfaces:**
```typescript
interface SessionManager {
  createSession(): SessionInfo
  getSession(sessionId: string): SessionInfo
  deleteSession(sessionId: string): void
  cleanupExpiredSessions(): void
}

interface SessionInfo {
  sessionId: string
  createdAt: Date
  expiresAt: Date
  language: LanguageCode
  conversationHistory: ConversationExchange[]
}
```

## Data Models

### Core Data Structures

```typescript
// Patient Request Models
interface DocumentPayload {
  file: File
  fileType: 'pdf' | 'image' | 'text'
  fileName: string
}

interface QuestionPayload {
  question: string
  questionType: 'diagnosis' | 'medication' | 'lab_report' | 'follow_up' | 'general'
}

// Response Models
interface ExplanationResponse {
  content: string
  disclaimer: string
  confidence: 'high' | 'medium' | 'low'
  suggestDoctorConsultation: boolean
  relatedEntities: string[]
  timestamp: Date
}

interface ErrorMessage {
  code: string
  message: string
  userFriendlyMessage: string
  suggestedAction: string
}

// Processing Models
interface DocumentMetadata {
  originalFileName: string
  fileType: string
  fileSize: number
  uploadedAt: Date
  processedAt: Date
}

interface ConversationHistory {
  sessionId: string
  exchanges: ConversationExchange[]
  startedAt: Date
}

// Configuration Models
interface SystemConfiguration {
  maxDocumentSize: number
  ocrTimeout: number
  nerTimeout: number
  llmTimeout: number
  confidenceThreshold: number
  maxConversationHistory: number
  supportedLanguages: LanguageCode[]
}

type LanguageCode = 'en' | 'es' | 'fr' | 'hi' | 'zh' | 'ar' | 'pt'
```

### Data Flow

1. **Document Upload Flow:**
   ```
   Patient → UI → API Gateway → OCR Engine → NER Engine → 
   Entity Validator → Context Builder → Explanation Engine → 
   Guardrail System → API Gateway → UI → Patient
   ```

2. **Question Flow (No Document):**
   ```
   Patient → UI → API Gateway → Context Builder (Minimal Mode) → 
   Explanation Engine → Guardrail System → API Gateway → UI → Patient
   ```

3. **Red Flag Detection Flow:**
   ```
   Patient Symptoms → Red Flag Detector → Urgency Determination → 
   Triage Recommendation → Guardrail System → Patient
   ```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*


### Document Processing Properties

**Property 1: Document Type Routing**
*For any* uploaded document (PDF, image, or text), the system should route it to the appropriate processing pipeline: PDFs and images to OCR extraction, text documents to direct processing.
**Validates: Requirements 1.1, 1.2, 1.3**

**Property 2: Processing Pipeline Sequence**
*For any* document, after text extraction completes, NER entity extraction must be invoked, followed by entity validation.
**Validates: Requirements 1.4, 1.5**

**Property 3: Extraction Failure Handling**
*For any* document where text extraction fails, the system should notify the patient and request resubmission without proceeding to NER.
**Validates: Requirements 1.6**

**Property 4: Document Deletion After Processing**
*For any* uploaded medical document, after processing completes (successfully or with errors), the document must be deleted immediately and not exist in persistent storage.
**Validates: Requirements 1.7, 1.8, 10.3, 10.4**

### Explanation Generation Properties

**Property 5: Simple Language Generation**
*For any* generated explanation, the text should meet readability standards for non-medical audiences (e.g., Flesch-Kincaid grade level ≤ 8).
**Validates: Requirements 2.1, 15.3**

**Property 6: Entity Grounding**
*For any* generated explanation containing medical information, all medical terms and facts must be traceable to extracted medical entities or knowledge base sources.
**Validates: Requirements 2.2, 2.7, 17.1**

**Property 7: Abstention on Insufficient Context**
*For any* patient question where extracted entities are insufficient or ambiguous, the system should abstain from answering and recommend consulting a healthcare professional.
**Validates: Requirements 2.3, 3.7, 5.6, 6.5, 7.6, 9.7, 14.1, 14.3**

**Property 8: Disclaimer Inclusion**
*For any* generated response, a visible disclaimer must be present indicating the system is not a substitute for professional medical advice.
**Validates: Requirements 2.6, 15.4**

### Safety Guardrail Properties

**Property 9: Prohibited Content Blocking**
*For any* response generated by the Explanation Engine, the Guardrail System must block content that:
- Diagnoses diseases or conditions
- Prescribes medications
- Recommends treatment alternatives
- Contradicts doctor decisions
- Predicts disease risk from symptoms
- Suggests medication dosage changes

If any prohibited content is detected, the system must generate a safe alternative response.
**Validates: Requirements 2.4, 2.5, 3.4, 3.5, 3.6, 4.4, 7.4, 7.5, 9.1, 9.2, 9.3, 9.4, 9.6, 13.3**

**Property 10: Guardrail Validation for All Responses**
*For any* response generated by the Explanation Engine, the Guardrail System must validate it against safety rules before delivery to the patient.
**Validates: Requirements 9.5**

**Property 11: Low Confidence Abstention**
*For any* response where the Explanation Engine or Guardrail System has confidence below the safety threshold, the system must abstain and recommend professional consultation.
**Validates: Requirements 9.7, 14.1**

### Medication Guidance Properties

**Property 12: Medication Purpose Explanation**
*For any* prescribed medication in extracted entities, explanations should include the medication's purpose based on prescription details or knowledge base information.
**Validates: Requirements 3.1**

**Property 13: Usage Instructions from Prescription**
*For any* medication with dosage and frequency information in the prescription, the explanation should include usage instructions that match the prescription details.
**Validates: Requirements 3.2**

**Property 14: Side Effects from Guidelines**
*For any* medication explanation, common side effects should be sourced from the public medical guidelines knowledge base, not generated.
**Validates: Requirements 3.3**

### Lab Report Interpretation Properties

**Property 15: Lab Value Simplification**
*For any* lab report with extracted values, the interpretation should explain each test in simple language meeting readability standards.
**Validates: Requirements 4.1**

**Property 16: Abnormal Value Highlighting**
*For any* lab value outside the normal range specified in the report, the system should explicitly highlight it as abnormal in the interpretation.
**Validates: Requirements 4.2**

**Property 17: Lab Test Purpose Explanation**
*For any* lab test in the report, the interpretation should explain what the test measures.
**Validates: Requirements 4.3**

**Property 18: Urgent Lab Value Escalation**
*For any* lab value indicating potential medical urgency (based on severity thresholds), the system should recommend consulting the doctor.
**Validates: Requirements 4.5**

**Property 19: Minimal Input Mode Activation**
*For any* patient interaction where no medical documents are uploaded, the system should operate in Minimal Input Mode.
**Validates: Requirements 4.6, 8.1**

### Follow-Up Guidance Properties

**Property 20: Care Roadmap from Documents**
*For any* follow-up guidance request with medical documents, the generated roadmap should only include information present in the extracted entities.
**Validates: Requirements 5.1**

**Property 21: Medication Schedule Inclusion**
*For any* prescription with schedule information (frequency, duration), the follow-up guidance should include the medication schedule.
**Validates: Requirements 5.2**

**Property 22: Appointment Reminder Inclusion**
*For any* medical document containing follow-up appointment information, the guidance should include appointment reminders.
**Validates: Requirements 5.3**

**Property 23: Diagnosis-Relevant Lifestyle Advice**
*For any* follow-up guidance with a known diagnosis, lifestyle advice should be relevant to the diagnosed condition based on knowledge base guidelines.
**Validates: Requirements 5.4**

**Property 24: Follow-Up Instruction Preservation**
*For any* doctor-specified follow-up instructions in medical documents, the system should present them without modification.
**Validates: Requirements 5.5**

### Preventive Health Properties

**Property 25: Guideline-Based Preventive Advice**
*For any* preventive advice request, recommendations should be sourced from public health guidelines in the knowledge base.
**Validates: Requirements 6.1**

**Property 26: Condition-Tailored Preventive Advice**
*For any* preventive advice request with a known diagnosis, recommendations should be tailored to the patient's condition.
**Validates: Requirements 6.2**

**Property 27: Minimal Mode Preventive Advice**
*For any* preventive advice request in Minimal Input Mode (no documents), the system should provide general preventive recommendations.
**Validates: Requirements 6.3**

**Property 28: Preventive Advice Non-Contradiction**
*For any* preventive advice generated, it must not contradict doctor instructions present in medical documents.
**Validates: Requirements 6.4**

### Red Flag Detection Properties

**Property 29: Symptom Pattern Matching**
*For any* patient-described symptoms, the system should check them against rule-based red flag symptom patterns in the knowledge base.
**Validates: Requirements 7.1**

**Property 30: Red Flag Consultation Recommendation**
*For any* symptom description matching red flag patterns, the system should recommend immediate medical consultation.
**Validates: Requirements 7.2**

**Property 31: Emergency Service Recommendation**
*For any* symptom description matching emergency-level red flag patterns, the system should recommend emergency services (911/emergency hotline).
**Validates: Requirements 7.3**

**Property 32: Uncertain Symptom Escalation**
*For any* symptom description where severity is uncertain, the system should recommend consulting a healthcare professional.
**Validates: Requirements 7.6**

### Minimal Input Mode Properties

**Property 33: General Health Education in Minimal Mode**
*For any* health-related question in Minimal Input Mode, the system should provide general health education from the knowledge base.
**Validates: Requirements 8.2**

**Property 34: Common Medication Explanation in Minimal Mode**
*For any* medication question in Minimal Input Mode (without prescription documents), the system should provide general medication usage information from the knowledge base.
**Validates: Requirements 8.3**

**Property 35: Professional Consultation Encouragement in Minimal Mode**
*For any* specific medical question in Minimal Input Mode, the system should encourage consulting healthcare professionals.
**Validates: Requirements 8.5**

**Property 36: No Personalized Advice in Minimal Mode**
*For any* response in Minimal Input Mode, the system must not provide personalized medical advice specific to the patient's condition.
**Validates: Requirements 8.6**

### Privacy and Security Properties

**Property 37: Comprehensive Data Encryption**
*For any* medical document or PHI, the system should encrypt data in transit (during upload/download) and at rest (during processing).
**Validates: Requirements 10.1, 10.2, 10.6**

**Property 38: Session Isolation**
*For any* two concurrent patient sessions, data processed in one session must not be accessible to the other session.
**Validates: Requirements 10.7**

**Property 39: Conversation History Deletion**
*For any* patient session, after session termination, conversation history must be deleted and not retained.
**Validates: Requirements 12.6**

**Property 40: PHI-Free Audit Logs**
*For any* audit log entry (guardrail interventions, errors, interactions), the log must not contain Protected Health Information.
**Validates: Requirements 20.3**

### Multi-Language Properties

**Property 41: Response Language Matching**
*For any* patient interaction with a selected language, generated explanations should be in the selected language.
**Validates: Requirements 11.1**

**Property 42: Multi-Language Document Processing**
*For any* medical document in a language different from the interface language, the system should attempt to process the document in its original language.
**Validates: Requirements 11.3**

**Property 43: Translation Quality Abstention**
*For any* translation with low confidence or quality, the system should notify the patient and recommend professional translation services.
**Validates: Requirements 11.4**

### Conversational Interface Properties

**Property 44: Natural Language Processing**
*For any* text question submitted by a patient, the system should process it as natural language input.
**Validates: Requirements 12.1**

**Property 45: Speech-to-Text Conversion**
*For any* voice question submitted by a patient (where voice input is supported), the system should convert speech to text before processing.
**Validates: Requirements 12.2**

**Property 46: Conversation Context Maintenance**
*For any* follow-up question within a session, the system should maintain and reference previous exchanges in the conversation history.
**Validates: Requirements 12.3, 12.4**

**Property 47: Context Loss Clarification**
*For any* question where conversation context is insufficient to understand the question, the system should request clarification from the patient.
**Validates: Requirements 12.5**

### Trust and Authority Properties

**Property 48: Doctor Authority Framing**
*For any* explanation of medical decisions, the system should present doctor decisions as authoritative and attribute them to the doctor, framing responses as explanations rather than medical advice.
**Validates: Requirements 13.1, 13.2, 13.5**

**Property 49: Disagreement Escalation**
*For any* patient question suggesting disagreement with doctor decisions, the system should encourage direct doctor consultation rather than providing alternative perspectives.
**Validates: Requirements 13.4**

### Uncertainty Handling Properties

**Property 50: Confidence Indicator Inclusion**
*For any* response where confidence is not high, the system should include confidence indicators to signal uncertainty.
**Validates: Requirements 14.4**

**Property 51: Limitation Acknowledgment**
*For any* response based on incomplete information, the system should explicitly acknowledge its limitations.
**Validates: Requirements 14.5**

### Error Handling Properties

**Property 52: User-Friendly Error Messages**
*For any* error condition, the system should display user-friendly error messages that explain the issue and suggest corrective actions.
**Validates: Requirements 15.6**

**Property 53: OCR Failure Handling**
*For any* OCR processing failure, the system should notify the patient and suggest document quality improvements (better lighting, higher resolution, etc.).
**Validates: Requirements 18.1**

**Property 54: NER Failure Fallback**
*For any* NER extraction failure, the system should attempt text-based processing and notify the patient of limitations.
**Validates: Requirements 18.2**

**Property 55: Engine Unavailability Handling**
*For any* Explanation Engine unavailability, the system should display a maintenance message and suggest retry timing.
**Validates: Requirements 18.3**

**Property 56: Network Failure Recovery**
*For any* network connectivity loss during patient input, the system should preserve the input and retry when connectivity is restored.
**Validates: Requirements 18.4**

**Property 57: PHI-Free Error Logging**
*For any* critical error, the system should log error details for administrators without exposing Protected Health Information.
**Validates: Requirements 18.5**

### Performance Properties

**Property 58: Response Time Thresholds**
*For any* patient interaction, the system should meet response time requirements:
- Text questions without documents: ≤ 5 seconds
- OCR processing for standard documents: ≤ 30 seconds
- Explanation generation after processing: ≤ 10 seconds
**Validates: Requirements 16.1, 16.2, 16.3**

**Property 59: High Load Queuing**
*For any* request received when system load is high, the system should queue the request and notify the patient of expected wait times.
**Validates: Requirements 16.4**

### Audit and Compliance Properties

**Property 60: Guardrail Intervention Logging**
*For any* guardrail intervention (blocking, modification, abstention), the system should create an audit log entry.
**Validates: Requirements 20.1**

**Property 61: Error and Failure Logging**
*For any* system error or failure, the system should create an audit log entry with error details.
**Validates: Requirements 20.2**

**Property 62: Anonymized Session Logging**
*For any* logged interaction, the system should use anonymized session identifiers rather than patient-identifying information.
**Validates: Requirements 20.4**

**Property 63: Audit Log Retention**
*For any* audit log entry, the system should retain it for the compliance-required review period.
**Validates: Requirements 20.5**

**Property 64: Audit Log Export**
*For any* compliance reporting request, the system should support exporting audit logs in a standard format.
**Validates: Requirements 20.6**

**Property 65: Source Citation**
*For any* medical information provided from the knowledge base, the system should cite sources when appropriate.
**Validates: Requirements 17.5**

## Error Handling

### Error Categories

1. **Document Processing Errors**
   - OCR extraction failures (poor image quality, unsupported format)
   - NER extraction failures (unrecognized medical terminology)
   - Entity validation failures (inconsistent or incomplete entities)

2. **AI Generation Errors**
   - LLM timeout or unavailability
   - Guardrail violations (prohibited content detected)
   - Low confidence responses (below safety threshold)

3. **System Errors**
   - Network connectivity issues
   - Service unavailability
   - Session management failures

4. **User Input Errors**
   - Invalid file formats
   - File size exceeding limits
   - Malformed questions or inputs

### Error Handling Strategies

**Graceful Degradation:**
- When OCR fails, allow manual text input
- When NER fails, attempt text-based processing
- When documents unavailable, operate in Minimal Input Mode

**Safe Defaults:**
- When uncertain, abstain and recommend professional consultation
- When guardrails detect violations, block and provide safe alternative
- When confidence is low, include explicit uncertainty indicators

**User Communication:**
- Provide clear, actionable error messages
- Suggest corrective actions (improve document quality, rephrase question)
- Indicate expected resolution time for system errors

**Logging and Monitoring:**
- Log all errors for system administrators (without PHI)
- Track guardrail interventions for safety monitoring
- Monitor performance metrics for SLA compliance

### Error Recovery

**Retry Logic:**
- Automatic retry for transient network failures
- Preserve user input during connectivity loss
- Exponential backoff for service unavailability

**Fallback Mechanisms:**
- Text-based processing when NER fails
- Minimal Input Mode when documents unavailable
- Knowledge base lookup when entity extraction incomplete

**Session Management:**
- Maintain session state during temporary failures
- Clean up sessions after terminal errors
- Prevent data leakage across failed sessions

## Testing Strategy

### Dual Testing Approach

The testing strategy employs both unit testing and property-based testing as complementary approaches:

**Unit Tests:**
- Verify specific examples and edge cases
- Test integration points between components
- Validate error conditions and failure modes
- Test specific document formats and content types
- Verify UI rendering and user interactions

**Property-Based Tests:**
- Verify universal properties across all inputs
- Test safety guardrails with randomized content
- Validate privacy properties with various data types
- Test language support across multiple languages
- Verify performance under varied load conditions

### Property-Based Testing Configuration

**Testing Library:** Use a property-based testing library appropriate for the implementation language:
- Python: Hypothesis
- TypeScript/JavaScript: fast-check
- Java: jqwik
- Go: gopter

**Test Configuration:**
- Minimum 100 iterations per property test (due to randomization)
- Each property test must reference its design document property
- Tag format: **Feature: medipath-ai, Property {number}: {property_text}**

**Example Property Test Structure:**
```python
# Feature: medipath-ai, Property 9: Prohibited Content Blocking
@given(st.text(), st.sampled_from(['diagnosis', 'prescription', 'treatment_alternative']))
def test_guardrail_blocks_prohibited_content(response_text, violation_type):
    # Inject prohibited content into response
    response = inject_violation(response_text, violation_type)
    
    # Run through guardrail
    result = guardrail_system.validate_response(response, context)
    
    # Assert violation detected and blocked
    assert result.approved == False
    assert any(v.rule.value == f'no_{violation_type}' for v in result.violations)
    assert result.action == 'abstain' or result.action == 'modify'
```

### Testing Coverage by Component

**OCR Engine:**
- Unit tests: Specific PDF/image formats, various document qualities
- Property tests: All documents produce text or error, no crashes on malformed input

**NER Engine:**
- Unit tests: Known medical documents with expected entities
- Property tests: All extracted entities have valid types, entity extraction doesn't modify original text

**Guardrail System:**
- Unit tests: Specific prohibited phrases and patterns
- Property tests: All prohibited content types blocked, all safe content approved

**Explanation Engine:**
- Unit tests: Specific medical scenarios with expected explanations
- Property tests: All responses grounded in entities, all responses include disclaimers

**Privacy/Security:**
- Unit tests: Specific encryption algorithms and key management
- Property tests: All PHI encrypted, all documents deleted after processing, no PHI in logs

**Multi-Language:**
- Unit tests: Specific language pairs and translations
- Property tests: Response language matches selected language for all supported languages

### Integration Testing

**End-to-End Flows:**
1. Document upload → OCR → NER → Explanation → Response
2. Question without document → Minimal Mode → Knowledge Base → Response
3. Red flag symptoms → Detection → Triage → Urgent recommendation

**Component Integration:**
- API Gateway ↔ Processing Pipeline
- Processing Pipeline ↔ AI Layer
- AI Layer ↔ Knowledge Layer
- All layers ↔ Security Layer

### Safety Testing

**Adversarial Testing:**
- Attempt to elicit diagnoses through various phrasings
- Attempt to get prescription recommendations
- Attempt to get contradictions of doctor decisions
- Test with medical misinformation

**Boundary Testing:**
- Test confidence thresholds for abstention
- Test entity completeness thresholds
- Test red flag detection sensitivity
- Test performance under load

### Privacy Testing

**Data Lifecycle Testing:**
- Verify document deletion after processing
- Verify session cleanup after termination
- Verify no PHI in logs or persistent storage
- Verify encryption at all stages

**Session Isolation Testing:**
- Verify no data leakage between sessions
- Verify session-specific context maintenance
- Verify proper session termination

### Performance Testing

**Load Testing:**
- Test response times under normal load
- Test queuing behavior under high load
- Test resource utilization and scaling

**Stress Testing:**
- Test with maximum document sizes
- Test with complex medical documents
- Test with high concurrency

### Compliance Testing

**Audit Log Testing:**
- Verify all required events logged
- Verify no PHI in audit logs
- Verify log retention and export

**Regulatory Compliance:**
- HIPAA compliance verification (if applicable)
- GDPR compliance verification (if applicable)
- Accessibility compliance (WCAG guidelines)

## Implementation Notes

### Technology Recommendations

**OCR Engine:**
- Tesseract OCR for open-source option
- Google Cloud Vision API or AWS Textract for cloud-based option
- Consider multi-language OCR support

**NER Engine:**
- spaCy with medical models (scispaCy)
- BioBERT or ClinicalBERT for medical entity recognition
- Custom fine-tuned models for specific medical domains

**LLM for Explanation Engine:**
- GPT-4 or Claude for high-quality explanations
- Open-source alternatives: Llama 2, Mistral
- Consider medical-domain fine-tuned models

**Guardrail Implementation:**
- Rule-based pattern matching for prohibited content
- Classifier models for content categorization
- Confidence scoring for abstention decisions
- Multiple layers: pre-generation prompts + post-generation validation

**Knowledge Base:**
- Vector database for medical guidelines (Pinecone, Weaviate)
- Structured database for medications and red flags (PostgreSQL)
- Regular updates from public health sources (CDC, WHO, NIH)

### Deployment Considerations

**Infrastructure:**
- Containerized services (Docker/Kubernetes)
- Serverless functions for API Gateway (AWS Lambda, Google Cloud Functions)
- GPU instances for LLM inference
- CDN for frontend delivery

**Security:**
- TLS 1.3 for all communications
- AES-256 for data at rest
- Key management service (AWS KMS, Google Cloud KMS)
- Network isolation for processing pipeline

**Scalability:**
- Horizontal scaling for API Gateway and processing pipeline
- Load balancing across LLM instances
- Caching for knowledge base queries
- Queue-based request handling for high load

**Monitoring:**
- Real-time monitoring of response times
- Guardrail intervention tracking
- Error rate monitoring
- Audit log analysis

### Development Phases

**Phase 1: Core Pipeline**
- Document upload and OCR
- Basic NER extraction
- Simple explanation generation
- Basic guardrails

**Phase 2: Safety and Privacy**
- Comprehensive guardrail system
- Privacy controls and encryption
- Session management
- Audit logging

**Phase 3: Enhanced Features**
- Multi-language support
- Voice input
- Red flag detection
- Minimal Input Mode

**Phase 4: Optimization**
- Performance optimization
- Advanced error handling
- Knowledge base expansion
- Accessibility features

### Maintenance and Updates

**Knowledge Base Updates:**
- Regular updates from public health sources
- Version control for medical guidelines
- A/B testing for new content

**Model Updates:**
- Periodic retraining of NER models
- LLM version upgrades with safety testing
- Guardrail rule refinement based on interventions

**Monitoring and Improvement:**
- User feedback collection (anonymized)
- Guardrail intervention analysis
- Performance metric tracking
- Continuous safety audits
