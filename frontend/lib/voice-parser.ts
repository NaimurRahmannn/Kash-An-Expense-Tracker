import { formatDateForAPI } from "@/lib/date";
import type { ExpenseCategory } from "@/types/expense";

export type ParsedVoiceExpense = {
  title: string;
  amount: string;
  category: ExpenseCategory | "";
  expense_date: string;
  note: string;
  confidence: "high" | "low";
  transcript: string;
  error?: string;
};

const NO_SPEECH_ERROR = "No speech detected";

const VOICE_DATE_ERROR =
  "Date is missing or invalid. Use today, yesterday, or YYYY-MM-DD";

const CATEGORY_ERROR =
  "Category not recognized. Say food, transport, housing, entertainment, shopping, healthcare, education, utilities, or other.";

const CATEGORY_MAP: Record<string, ExpenseCategory> = {
  food: "Food",
  transport: "Transport",
  transportation: "Transport",
  housing: "Housing",
  entertainment: "Entertainment",
  shopping: "Shopping",
  healthcare: "Healthcare",
  "health care": "Healthcare",
  education: "Education",
  utilities: "Utilities",
  utility: "Utilities",
  other: "Other",
};

const categoryFillers = new Set(["on", "for", "at"]);

function createLowConfidenceResult(
  transcript: string,
  error: string,
  overrides: Partial<ParsedVoiceExpense> = {},
): ParsedVoiceExpense {
  return {
    title: "",
    amount: "",
    category: "",
    expense_date: "",
    note: transcript ? `Voice note: ${transcript}` : "",
    confidence: "low",
    transcript,
    error,
    ...overrides,
  };
}

function normalizeTranscript(transcript: string) {
  return transcript
    .toLowerCase()
    .replace(/(\d),(?=\d)/g, "$1")
    .replace(/\s+/g, " ")
    .trim();
}

function cleanToken(token: string) {
  return token.replace(/^[^a-z0-9-]+|[^a-z0-9-]+$/gi, "");
}

function isNumericToken(token: string) {
  return /^-?\d+(?:\.\d+)?$/.test(cleanToken(token));
}

function formatTitle(title: string) {
  const trimmedTitle = title.trim();

  if (!trimmedTitle) {
    return "";
  }

  return `${trimmedTitle.charAt(0).toUpperCase()}${trimmedTitle.slice(1)}`;
}

function cleanTitle(titleWords: string[]) {
  let title = titleWords.map(cleanToken).filter(Boolean).join(" ").trim();

  if (title.startsWith("add expense ")) {
    title = title.slice("add expense ".length).trim();
  } else if (title === "add expense") {
    title = "";
  } else if (title.startsWith("add ")) {
    title = title.slice("add ".length).trim();
  } else if (title === "add") {
    title = "";
  }

  return formatTitle(title);
}

function parseStrictDate(dateToken: string) {
  const cleanedToken = cleanToken(dateToken);

  if (cleanedToken === "today") {
    return formatDateForAPI(new Date());
  }

  if (cleanedToken === "yesterday") {
    const yesterday = new Date();
    yesterday.setDate(yesterday.getDate() - 1);

    return formatDateForAPI(yesterday);
  }

  if (!/^\d{4}-\d{2}-\d{2}$/.test(cleanedToken)) {
    return null;
  }

  const [year, month, day] = cleanedToken.split("-").map(Number);
  const parsedDate = new Date(year, month - 1, day);
  const formattedDate = formatDateForAPI(parsedDate);

  return formattedDate === cleanedToken ? cleanedToken : null;
}

function findDateToken(words: string[]) {
  for (let index = words.length - 1; index >= 0; index -= 1) {
    const token = cleanToken(words[index]);

    if (token === "today" || token === "yesterday") {
      return {
        date: parseStrictDate(token),
        index,
      };
    }

    if (/^\d{4}-\d{2}-\d{2}$/.test(token)) {
      const parsedDate = parseStrictDate(token);

      return {
        date: parsedDate,
        index,
      };
    }
  }

  return null;
}

function parseCategory(categoryWords: string[]) {
  const categoryCandidate = categoryWords
    .map(cleanToken)
    .filter((word) => word && !categoryFillers.has(word))
    .join(" ")
    .replace(/\s+/g, " ")
    .trim();

  return categoryCandidate ? CATEGORY_MAP[categoryCandidate] : undefined;
}

export function parseExpenseFromTranscript(
  transcript: string,
): ParsedVoiceExpense {
  const originalTranscript = transcript.trim();

  if (!originalTranscript) {
    return createLowConfidenceResult("", NO_SPEECH_ERROR);
  }

  const normalizedTranscript = normalizeTranscript(originalTranscript);
  const words = normalizedTranscript.split(" ").map(cleanToken).filter(Boolean);
  const amountIndex = words.findIndex(isNumericToken);

  if (amountIndex === -1) {
    return createLowConfidenceResult(
      originalTranscript,
      "Amount not found in transcript",
    );
  }

  const amount = words[amountIndex];
  const amountValue = Number(amount);

  if (!Number.isFinite(amountValue) || amountValue <= 0) {
    return createLowConfidenceResult(originalTranscript, "Amount must be positive", {
      amount,
    });
  }

  const title = cleanTitle(words.slice(0, amountIndex));

  if (!title) {
    return createLowConfidenceResult(originalTranscript, "Title is missing", {
      amount,
    });
  }

  const remainingWords = words.slice(amountIndex + 1);
  const dateResult = findDateToken(remainingWords);

  if (!dateResult || !dateResult.date) {
    return createLowConfidenceResult(originalTranscript, VOICE_DATE_ERROR, {
      amount,
      title,
    });
  }

  const category = parseCategory(remainingWords.slice(0, dateResult.index));

  if (!category) {
    return createLowConfidenceResult(originalTranscript, CATEGORY_ERROR, {
      amount,
      expense_date: dateResult.date,
      title,
    });
  }

  return {
    title,
    amount,
    category,
    expense_date: dateResult.date,
    note: `Voice note: ${originalTranscript}`,
    confidence: "high",
    transcript: originalTranscript,
  };
}
