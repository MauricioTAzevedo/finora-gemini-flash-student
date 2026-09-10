/**
 * Formats integer minor units (cents) as Brazilian Real currency: R$ 1.234,56
 */
export function formatBRL(amountMinor: number): string {
  const isNegative = amountMinor < 0;
  const absCents = Math.abs(amountMinor);
  const whole = Math.floor(absCents / 100);
  const fraction = absCents % 100;

  const wholeFormatted = whole.toLocaleString('pt-BR');
  const fractionFormatted = fraction.toString().padStart(2, '0');

  const prefix = isNegative ? '-R$ ' : 'R$ ';
  return `${prefix}${wholeFormatted},${fractionFormatted}`;
}

/**
 * Formats ISO date string into Brazilian date format: DD/MM/YYYY
 */
export function formatDateBR(dateStr: string): string {
  if (!dateStr) return '';
  const d = new Date(dateStr);
  return d.toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  });
}
