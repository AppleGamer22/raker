import { CropIcon, ImageIcon } from "lucide-react";
import {
	useCallback,
	useEffect,
	useRef,
	useState,
	type CSSProperties,
	type ReactElement,
	type SyntheticEvent,
	type PointerEvent as ReactPointerEvent,
} from "react";
import { ResizableBox, type ResizeCallbackData } from "react-resizable";

import "react-resizable/css/styles.css";

import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { Button } from "@/components/ui/button";
import {
	Carousel,
	CarouselContent,
	CarouselItem,
	CarouselNext,
	CarouselPrevious,
	type CarouselApi,
} from "@/components/ui/carousel";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { GoogleMapsLink } from "@/components/ui/svgs/google-maps";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";

export function postTypeString(type: PostType): string {
	switch (type) {
		case PostType.Instagram:
			return "instagram";
		case PostType.Highlight:
			return "highlight";
		case PostType.Story:
			return "story";
		case PostType.TikTok:
			return "tiktok";
		case PostType.Snapchat:
			return "snapchat";
		case PostType.VSCO:
			return "vsco";
	}
}

export function FilesCarousel({
	username,
	post: { postType, postOwner, coordinates, post, files },
}: {
	username: string;
	post: ScrapeResponse;
}) {
	const [api, setApi] = useState<CarouselApi>();
	const [selectedIndex, setSelectedIndex] = useState(0);
	const prevHeightsRef = useRef<Map<number, number>>(new Map());

	const syncSelectedSlideHeight = useCallback((emblaApi: CarouselApi) => {
		if (!emblaApi) {
			return;
		}

		const selectedIndex = emblaApi.selectedScrollSnap();
		const selectedSlide = emblaApi.slideNodes()[selectedIndex];
		if (!selectedSlide) {
			return;
		}

		const h = selectedSlide.offsetHeight;
		const prev = prevHeightsRef.current.get(selectedIndex) || 0;

		if (h > prev) {
			prevHeightsRef.current.set(selectedIndex, h);
			emblaApi.containerNode().style.height = `${h}px`;
		} else if (prev > 0) {
			// Do not shrink below the previously observed max height for this slide.
			emblaApi.containerNode().style.height = `${prev}px`;
		} else {
			emblaApi.containerNode().style.height = `${h}px`;
		}
	}, []);

	useEffect(() => {
		if (!api) {
			return;
		}

		const onSelect = () => {
			setSelectedIndex(api.selectedScrollSnap());
			syncSelectedSlideHeight(api);
		};

		const onReInit = () => {
			setSelectedIndex(api.selectedScrollSnap());
			syncSelectedSlideHeight(api);
		};

		onSelect();
		api.on("reInit", onReInit);
		api.on("select", onSelect);
		api.on("settle", onSelect);

		return () => {
			api.off("reInit", onReInit);
			api.off("select", onSelect);
			api.off("settle", onSelect);
		};
	}, [api, syncSelectedSlideHeight]);

	return files.length > 1 ? (
		<Carousel opts={{ loop: true }} setApi={setApi}>
			<CarouselContent className="items-center">
				{files.map((file, i) => (
					<CarouselItem
						key={`file-${postType}-${postOwner}-${post}-${i}`}
						className="flex items-center justify-center self-center"
					>
						<FileDisplay
							username={username}
							post={{ postType, postOwner, coordinates } as ScrapeResponse}
							file={file}
							onMediaLoad={() => {
								api?.reInit();
								if (api) {
									requestAnimationFrame(() => syncSelectedSlideHeight(api));
								}
							}}
							withCrop
							withCoordinates
						/>
					</CarouselItem>
				))}
			</CarouselContent>
			<CarouselPrevious className="top-1/2 left-2 z-10 -translate-y-1/2" />
			<CarouselNext className="top-1/2 right-2 z-10 -translate-y-1/2" />
			<div className="flex items-center justify-center gap-2 pt-2" role="group" aria-label="Slide navigation">
				{files.map((_, index) => (
					<button
						type="button"
						key={`dot-${postType}-${postOwner}-${post}-${index}`}
						onClick={() => api?.scrollTo(index)}
						className={`h-2 rounded-full transition-all focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none ${
							index === selectedIndex
								? "w-6 bg-primary"
								: "w-2 bg-muted-foreground/40 hover:bg-muted-foreground/70"
						}`}
						aria-label={`Go to slide ${index + 1}`}
						aria-current={index === selectedIndex ? "true" : undefined}
					/>
				))}
			</div>
		</Carousel>
	) : (
		<FileDisplay
			username={username}
			post={{ postType, postOwner, coordinates } as ScrapeResponse}
			file={files[0]}
			withCrop
			withCoordinates
		/>
	);
}

export function FileDisplay({
	username,
	file,
	post: { postType, postOwner, coordinates },
	onMediaLoad,
	withCrop,
	withCoordinates,
	className = "h-auto w-full rounded-xl",
}: {
	username: string;
	file: string;
	post: ScrapeResponse;
	onMediaLoad?: () => void;
	withCrop?: boolean;
	withCoordinates?: boolean;
	className?: string;
}) {
	const url = `/api/storage/${username}/${postTypeString(postType)}/${postOwner}/${file}`;
	if (/\.(jpe?g)|(webp)|(heic)$/.test(file)) {
		const imgResult = <img src={url} onLoad={onMediaLoad} loading="lazy" className={className} />;
		return withCrop || withCoordinates ? (
			<div className="relative inline-block w-full rounded-xl">
				{imgResult}
				<div className="absolute top-2 left-2 z-10 flex flex-row gap-2">
					{withCrop && /\.(jpe?g)$/.test(file) && (
						<FileSheet
							file={file}
							post={{ postType, postOwner, coordinates } as ScrapeResponse}
							username={username}
							trigger={
								<Button
									variant="outline"
									size="icon"
									className="dark:bg-secondary dark:hover:bg-secondary/80"
								>
									<CropIcon />
								</Button>
							}
						/>
					)}
					{withCoordinates && postType === PostType.VSCO && coordinates && (
						<GoogleMapsLink coordinates={coordinates} size="icon" />
					)}
				</div>
			</div>
		) : (
			imgResult
		);
	} else if (/\.(mp4)|(webm)$/.test(file)) {
		return (
			<video
				src={url}
				onLoadedMetadata={onMediaLoad}
				preload="metadata"
				className={className}
				loop
				controls
				muted
			/>
		);
	} else {
		return <a href={url}>{url}</a>;
	}
}

export type CropRect = {
	x: number;
	y: number;
	width: number;
	height: number;
};

const CROP_HANDLE_SIZE = 10;
const MIN_CROP_SIZE = 40;

function clamp(value: number, min: number, max: number) {
	return Math.min(Math.max(value, min), max);
}

function clampRect(rect: CropRect, maxWidth: number, maxHeight: number): CropRect {
	const safeWidth = clamp(rect.width, MIN_CROP_SIZE, maxWidth);
	const safeHeight = clamp(rect.height, MIN_CROP_SIZE, maxHeight);
	const safeX = clamp(rect.x, 0, Math.max(0, maxWidth - safeWidth));
	const safeY = clamp(rect.y, 0, Math.max(0, maxHeight - safeHeight));
	return {
		x: safeX,
		y: safeY,
		width: safeWidth,
		height: safeHeight,
	};
}

function displayToNatural(rect: CropRect, scaleX: number, scaleY: number): CropRect {
	return {
		x: Math.round(rect.x * scaleX),
		y: Math.round(rect.y * scaleY),
		width: Math.round(rect.width * scaleX),
		height: Math.round(rect.height * scaleY),
	};
}

function naturalToDisplay(rect: CropRect, scaleX: number, scaleY: number): CropRect {
	return {
		x: rect.x / scaleX,
		y: rect.y / scaleY,
		width: rect.width / scaleX,
		height: rect.height / scaleY,
	};
}

function handleStyle(handle: string): CSSProperties {
	const base: CSSProperties = {
		position: "absolute",
		width: CROP_HANDLE_SIZE,
		height: CROP_HANDLE_SIZE,
	};

	switch (handle) {
		case "se":
			return { ...base, right: 0, bottom: 0 };
		case "sw":
			return { ...base, left: 0, bottom: 0 };
		case "ne":
			return { ...base, right: 0, top: 0 };
		case "nw":
			return { ...base, left: 0, top: 0 };
		case "n":
			return { ...base, top: 0, left: "50%", transform: "translate(-50%, 0)" };
		case "s":
			return { ...base, bottom: 0, left: "50%", transform: "translate(-50%, 0)" };
		case "e":
			return { ...base, right: 0, top: "50%", transform: "translate(0, -50%)" };
		case "w":
			return { ...base, left: 0, top: "50%", transform: "translate(0, -50%)" };
		default:
			return base;
	}
}

function CropPreview({
	username,
	file,
	post: { postType, postOwner },
	onCropChange,
	className = "h-full w-auto rounded-xl",
}: {
	username: string;
	file: string;
	post: ScrapeResponse;
	onCropChange?: (rect: CropRect | null) => void;
	className?: string;
}) {
	const imageRef = useRef<HTMLImageElement | null>(null);
	const dragStateRef = useRef<{
		startX: number;
		startY: number;
		originX: number;
		originY: number;
		width: number;
		height: number;
	} | null>(null);
	const resizeHandleRef = useRef<string | null>(null);
	const [activeHandle, setActiveHandle] = useState<string | null>(null);
	const [naturalSize, setNaturalSize] = useState({ width: 0, height: 0 });
	const [displaySize, setDisplaySize] = useState({ width: 0, height: 0 });
	const [cropNatural, setCropNatural] = useState<CropRect | null>(null);

	const updateDisplaySize = useCallback(() => {
		const img = imageRef.current;
		if (!img) {
			return;
		}
		const rect = img.getBoundingClientRect();
		const width = Math.round(rect.width);
		const height = Math.round(rect.height);
		setDisplaySize((prev) => (prev.width === width && prev.height === height ? prev : { width, height }));
	}, []);

	useEffect(() => {
		setCropNatural(null);
		setNaturalSize({ width: 0, height: 0 });
		setDisplaySize({ width: 0, height: 0 });
	}, [file]);

	useEffect(() => {
		const img = imageRef.current;
		if (!img || typeof ResizeObserver === "undefined") {
			return;
		}
		const observer = new ResizeObserver(() => updateDisplaySize());
		observer.observe(img);
		return () => observer.disconnect();
	}, [updateDisplaySize, file]);

	const handleImageLoad = useCallback(() => {
		const img = imageRef.current;
		if (!img) {
			return;
		}
		setNaturalSize({ width: img.naturalWidth, height: img.naturalHeight });
		updateDisplaySize();
	}, [updateDisplaySize]);

	const scaleX = naturalSize.width > 0 && displaySize.width > 0 ? naturalSize.width / displaySize.width : 1;
	const scaleY = naturalSize.height > 0 && displaySize.height > 0 ? naturalSize.height / displaySize.height : 1;

	useEffect(() => {
		if (cropNatural || displaySize.width === 0 || displaySize.height === 0) {
			return;
		}
		const initialDisplay = clampRect(
			{
				x: 0,
				y: 0,
				width: displaySize.width,
				height: displaySize.height,
			},
			displaySize.width,
			displaySize.height,
		);
		setCropNatural(displayToNatural(initialDisplay, scaleX, scaleY));
	}, [cropNatural, displaySize, scaleX, scaleY]);

	const cropDisplay = cropNatural
		? clampRect(naturalToDisplay(cropNatural, scaleX, scaleY), displaySize.width, displaySize.height)
		: null;
	const maxConstraints: [number, number] = cropDisplay
		? (() => {
				const handle = activeHandle ?? resizeHandleRef.current ?? "se";
				const maxWidth = handle.includes("w")
					? cropDisplay.x + cropDisplay.width
					: displaySize.width - cropDisplay.x;
				const maxHeight = handle.includes("n")
					? cropDisplay.y + cropDisplay.height
					: displaySize.height - cropDisplay.y;
				return [Math.max(MIN_CROP_SIZE, maxWidth), Math.max(MIN_CROP_SIZE, maxHeight)];
			})()
		: [MIN_CROP_SIZE, MIN_CROP_SIZE];

	const updateFromDisplay = useCallback(
		(displayRect: CropRect) => {
			if (displaySize.width === 0 || displaySize.height === 0) {
				return;
			}
			const clamped = clampRect(displayRect, displaySize.width, displaySize.height);
			setCropNatural(displayToNatural(clamped, scaleX, scaleY));
		},
		[displaySize, scaleX, scaleY],
	);

	const handleResize = useCallback(
		(_event: SyntheticEvent, data: ResizeCallbackData) => {
			if (!cropDisplay) {
				return;
			}
			const handle = data.handle ?? resizeHandleRef.current ?? "se";
			const nextWidth = data.size.width;
			const nextHeight = data.size.height;
			const shiftX = cropDisplay.width - nextWidth;
			const shiftY = cropDisplay.height - nextHeight;
			const nextX = handle.includes("w") ? cropDisplay.x + shiftX : cropDisplay.x;
			const nextY = handle.includes("n") ? cropDisplay.y + shiftY : cropDisplay.y;
			updateFromDisplay({
				x: nextX,
				y: nextY,
				width: nextWidth,
				height: nextHeight,
			});
		},
		[cropDisplay, updateFromDisplay],
	);

	const handleResizeStart = useCallback((_event: SyntheticEvent, data: ResizeCallbackData) => {
		const handle = data.handle ?? null;
		resizeHandleRef.current = handle;
		setActiveHandle(handle);
	}, []);

	const handleResizeStop = useCallback(() => {
		resizeHandleRef.current = null;
		setActiveHandle(null);
	}, []);

	const handleDragStart = useCallback(
		(event: ReactPointerEvent<HTMLDivElement>) => {
			if (!cropDisplay || event.button !== 0) {
				return;
			}
			event.preventDefault();
			dragStateRef.current = {
				startX: event.clientX,
				startY: event.clientY,
				originX: cropDisplay.x,
				originY: cropDisplay.y,
				width: cropDisplay.width,
				height: cropDisplay.height,
			};

			const handleMove = (moveEvent: PointerEvent) => {
				const dragState = dragStateRef.current;
				if (!dragState) {
					return;
				}
				const dx = moveEvent.clientX - dragState.startX;
				const dy = moveEvent.clientY - dragState.startY;
				updateFromDisplay({
					x: dragState.originX + dx,
					y: dragState.originY + dy,
					width: dragState.width,
					height: dragState.height,
				});
			};

			const handleUp = () => {
				dragStateRef.current = null;
				window.removeEventListener("pointermove", handleMove);
				window.removeEventListener("pointerup", handleUp);
			};

			window.addEventListener("pointermove", handleMove);
			window.addEventListener("pointerup", handleUp);
		},
		[cropDisplay, updateFromDisplay],
	);

	useEffect(() => {
		if (!onCropChange) {
			return;
		}
		if (!cropNatural || naturalSize.width === 0 || naturalSize.height === 0) {
			onCropChange(null);
			return;
		}
		onCropChange(cropNatural);
	}, [cropNatural, naturalSize, onCropChange]);

	const url = `/api/storage/${username}/${postTypeString(postType)}/${postOwner}/${file}`;

	return (
		<div className={cn("relative inline-block", className)}>
			<img ref={imageRef} src={url} onLoad={handleImageLoad} className={className} />
			{cropDisplay && (
				<ResizableBox
					width={cropDisplay.width}
					height={cropDisplay.height}
					className="crop-box absolute z-10 rounded-md border-2 border-primary/80"
					style={{ left: cropDisplay.x, top: cropDisplay.y, position: "absolute" }}
					handleSize={[CROP_HANDLE_SIZE, CROP_HANDLE_SIZE]}
					minConstraints={[MIN_CROP_SIZE, MIN_CROP_SIZE]}
					maxConstraints={maxConstraints}
					handle={(handle, ref) => (
						<span
							ref={ref}
							className={`react-resizable-handle react-resizable-handle-${handle} crop-handle`}
							style={handleStyle(handle)}
							onPointerDown={(event) => event.stopPropagation()}
						/>
					)}
					resizeHandles={["sw", "se", "nw", "ne", "w", "e", "n", "s"]}
					onResizeStart={handleResizeStart}
					onResize={handleResize}
					onResizeStop={handleResizeStop}
					lockAspectRatio={false}
					axis="both"
					transformScale={1}
				>
					<div className="h-full w-full cursor-move rounded-md bg-black/10" onPointerDown={handleDragStart} />
				</ResizableBox>
			)}
		</div>
	);
}

export function FileSheet({
	trigger,
	username,
	file,
	post,
}: {
	trigger: ReactElement;
	username: string;
	file: string;
	post: ScrapeResponse;
}) {
	return (
		<Sheet>
			<SheetTrigger render={trigger} />
			<SheetContent side="bottom" className="p-1 pt-2 data-[side=bottom]:h-[90vh]">
				<Tabs defaultValue="view" className="flex h-full flex-col items-center">
					<TabsList>
						<TabsTrigger value="view">
							<ImageIcon />
							View
						</TabsTrigger>
						<TabsTrigger value="crop">
							<CropIcon />
							Crop
						</TabsTrigger>
					</TabsList>

					<div className="flex min-h-0 flex-1 flex-col items-center justify-center">
						<div className="flex max-h-full min-h-0 flex-col items-center">
							<TabsContent value="view" className="max-h-full w-auto rounded-xl">
								<FileDisplay
									file={file}
									post={post}
									username={username}
									className="h-full w-auto rounded-xl"
								/>
							</TabsContent>
							<TabsContent value="crop" className="max-h-full w-auto rounded-xl">
								<CropPreview file={file} post={post} username={username} onCropChange={undefined} />
							</TabsContent>
						</div>
					</div>
				</Tabs>
			</SheetContent>
		</Sheet>
	);
}
