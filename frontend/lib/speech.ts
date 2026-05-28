export function isSpeechRecognitionSupported(): boolean {
  if (typeof window === "undefined") {
    return false;
  }

  return Boolean(window.SpeechRecognition || window.webkitSpeechRecognition);
}

export function createSpeechRecognition(): SpeechRecognition | null {
  if (typeof window === "undefined") {
    return null;
  }

  const Recognition =
    window.SpeechRecognition || window.webkitSpeechRecognition;

  if (!Recognition) {
    return null;
  }

  const recognition = new Recognition();
  recognition.lang = "en-US";
  recognition.interimResults = false;
  recognition.maxAlternatives = 1;
  recognition.continuous = false;

  return recognition;
}
