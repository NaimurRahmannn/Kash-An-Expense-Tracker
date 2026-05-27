function padDatePart(value: number) {
  return String(value).padStart(2, "0");
}

export function formatDateForAPI(date: Date): string {
  const year = date.getFullYear();
  const month = padDatePart(date.getMonth() + 1);
  const day = padDatePart(date.getDate());

  return `${year}-${month}-${day}`;
}

export function getCurrentMonthRange(): { dateFrom: string; dateTo: string } {
  const today = new Date();
  const firstDay = new Date(today.getFullYear(), today.getMonth(), 1);
  const lastDay = new Date(today.getFullYear(), today.getMonth() + 1, 0);

  return {
    dateFrom: formatDateForAPI(firstDay),
    dateTo: formatDateForAPI(lastDay),
  };
}

export function formatDisplayDate(dateString: string): string {
  const date = new Date(`${dateString}T00:00:00`);

  return new Intl.DateTimeFormat("en", {
    month: "short",
    day: "numeric",
    year: "numeric",
  }).format(date);
}
