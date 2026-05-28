export {};

declare global {
  interface Window {
    SpeechRecognition?: SpeechRecognitionConstructor;
    webkitSpeechRecognition?: SpeechRecognitionConstructor;
  }

  interface SpeechRecognitionConstructor {
    new (): SpeechRecognition;
  }

  interface SpeechRecognition extends EventTarget {
    continuous: boolean;
    interimResults: boolean;
    lang: string;
    maxAlternatives: number;
    onend: ((this: SpeechRecognition, event: Event) => void) | null;
    onerror:
      | ((this: SpeechRecognition, event: SpeechRecognitionErrorEvent) => void)
      | null;
    onresult:
      | ((this: SpeechRecognition, event: SpeechRecognitionEvent) => void)
      | null;
    abort(): void;
    start(): void;
    stop(): void;
  }

  interface SpeechRecognitionAlternative {
    readonly confidence: number;
    readonly transcript: string;
  }

  interface SpeechRecognitionErrorEvent extends Event {
    readonly error: string;
    readonly message: string;
  }

  interface SpeechRecognitionEvent extends Event {
    readonly results: SpeechRecognitionResultList;
  }

  interface SpeechRecognitionResult {
    readonly isFinal: boolean;
    readonly length: number;
    [index: number]: SpeechRecognitionAlternative;
    item(index: number): SpeechRecognitionAlternative;
  }

  interface SpeechRecognitionResultList {
    readonly length: number;
    [index: number]: SpeechRecognitionResult;
    item(index: number): SpeechRecognitionResult;
  }
}
