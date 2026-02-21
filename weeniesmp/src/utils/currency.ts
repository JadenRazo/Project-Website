/**
 * Currency symbol mapping for common currencies
 */
const CURRENCY_SYMBOLS: Record<string, string> = {
  'USD': '$',
  'EUR': '€',
  'GBP': '£',
  'CAD': 'CA$',
  'AUD': 'A$',
  'JPY': '¥',
  'CNY': '¥',
  'KRW': '₩',
  'INR': '₹',
  'BRL': 'R$',
  'MXN': 'MX$',
  'RUB': '₽',
  'CHF': 'CHF',
  'SEK': 'kr',
  'NOK': 'kr',
  'DKK': 'kr',
  'PLN': 'zł',
  'TRY': '₺',
  'NZD': 'NZ$',
  'SGD': 'S$',
  'HKD': 'HK$',
  'ZAR': 'R',
}

/**
 * Converts a currency code (USD, EUR, etc.) to its symbol ($, €, etc.)
 */
export function formatCurrency(code: string | undefined): string {
  return CURRENCY_SYMBOLS[code ?? 'USD'] ?? code ?? '$'
}

/**
 * Formats a price with the currency symbol
 */
export function formatPrice(amount: number, currencyCode: string | undefined): string {
  const symbol = formatCurrency(currencyCode)
  return `${symbol}${amount.toFixed(2)}`
}
