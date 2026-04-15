import React from 'react';
import { Image as ImageIcon } from 'lucide-react';
import { AppBarContext } from '@/contexts/AppBarContext';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';

type ImagesResponse = {
  images: string[];
};

type Props = {
  logText?: string;
};

function extractImagePathFromLine(line: string): string | null {
  const match = line.match(/(\/[^\s"'<>|]+\.(?:jpg|jpeg|png|gif))/i);
  return match?.[1] ?? null;
}

export function DAGRunImageGallery({ logText = '' }: Props) {
  const appBarContext = React.useContext(AppBarContext);
  const [images, setImages] = React.useState<string[]>([]);
  const [selectedImage, setSelectedImage] = React.useState<string | null>(null);
  const [highlightedImage, setHighlightedImage] = React.useState<string | null>(
    null
  );

  const logLines = React.useMemo(() => {
    return logText
      .split('\n')
      .map((line) => line.trim())
      .filter(Boolean)
      .filter((line) => extractImagePathFromLine(line) !== null);
  }, [logText]);

  React.useEffect(() => {
    const controller = new AbortController();
    const params = new URLSearchParams({
      remoteNode: appBarContext.selectedRemoteNode || 'local',
    });
    if (logText.trim()) {
      params.set('log', logText);
    }

    fetch(`/api/v1/images?${params.toString()}`, { signal: controller.signal })
      .then(async (res) => {
        if (!res.ok) {
          throw new Error(`failed to load images: ${res.status}`);
        }
        return (await res.json()) as ImagesResponse;
      })
      .then((resp) => {
        setImages(resp.images || []);
      })
      .catch(() => {
        setImages([]);
      });

    return () => controller.abort();
  }, [appBarContext.selectedRemoteNode, logText]);

  const onClickLogLine = (line: string) => {
    const foundPath = extractImagePathFromLine(line);
    if (!foundPath) return;
    const matched = images.find((img) => img.endsWith(foundPath.split('/').pop() || ''));
    if (matched) {
      setHighlightedImage(matched);
      setSelectedImage(matched);
    }
  };

  if (images.length === 0) {
    return <div className="text-xs text-muted-foreground">No images found in logs.</div>;
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center gap-2">
        <ImageIcon className="h-4 w-4 text-muted-foreground" />
        <span className="text-sm font-medium">Generated Images ({images.length})</span>
      </div>

      {logLines.length > 0 && (
        <div className="rounded border border-border p-2 text-xs space-y-1 max-h-28 overflow-y-auto">
          {logLines.map((line, idx) => (
            <button
              key={`${line}-${idx}`}
              type="button"
              onClick={() => onClickLogLine(line)}
              className="block w-full text-left font-mono text-muted-foreground hover:text-foreground whitespace-normal break-all"
            >
              {line}
            </button>
          ))}
        </div>
      )}

      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-2">
        {images.map((src) => (
          <button
            key={src}
            type="button"
            onClick={() => setSelectedImage(src)}
            className={`rounded overflow-hidden border ${
              highlightedImage === src ? 'border-primary ring-1 ring-primary' : 'border-border'
            }`}
          >
            <img
              src={src}
              alt={src}
              loading="lazy"
              className="h-24 w-full object-cover bg-muted"
            />
          </button>
        ))}
      </div>

      <Dialog open={!!selectedImage} onOpenChange={(open) => !open && setSelectedImage(null)}>
        <DialogContent className="max-w-4xl p-3">
          <DialogHeader className="pb-1">
            <DialogTitle className="text-sm font-mono break-all">
              {selectedImage}
            </DialogTitle>
          </DialogHeader>
          {selectedImage && (
            <img
              src={selectedImage}
              alt={selectedImage}
              className="w-full max-h-[75vh] object-contain rounded"
            />
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}

