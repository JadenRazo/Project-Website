import { useSiteContent } from '../../lib/site-content';

/** React keeps ownership of this text during preview and later interactions. */
export default function SiteText({ entry, name, before, after }: { entry: string; name: string; before?: boolean; after?: boolean }) {
  const document = useSiteContent() as unknown as Record<string, Record<string, unknown>>;
  const value = document[entry]?.[name];
  return <>{before && ' '}<span data-rh={`${entry}.${name}`} data-rh-controlled style={{ display: 'contents' }}>{typeof value === 'string' ? value : ''}</span>{after && ' '}</>;
}
