import { useCallback, useEffect, useState } from "react";

import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import {
	Carousel,
	CarouselContent,
	CarouselItem,
	CarouselNext,
	CarouselPrevious,
	type CarouselApi,
} from "@/components/ui/carousel";

function postTypeString(type: PostType) {
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
	post: { postType, postOwner, post, files },
}: {
	username: string;
	post: ScrapeResponse;
}) {
	const [api, setApi] = useState<CarouselApi>();
	const [selectedIndex, setSelectedIndex] = useState(0);

	const syncSelectedSlideHeight = useCallback((emblaApi: CarouselApi) => {
		if (!emblaApi) {
			return;
		}

		const selectedIndex = emblaApi.selectedScrollSnap();
		const selectedSlide = emblaApi.slideNodes()[selectedIndex];
		if (!selectedSlide) {
			return;
		}

		emblaApi.containerNode().style.height = `${selectedSlide.offsetHeight}px`;
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
			<CarouselContent className="items-start">
				{files.map((file, i) => (
					<CarouselItem key={`file-${postType}-${postOwner}-${post}-${i}`} className="self-start">
						<FileDisplay
							username={username}
							post={{ postType, postOwner } as ScrapeResponse}
							file={file}
							onMediaLoad={() => {
								api?.reInit();
								if (api) {
									requestAnimationFrame(() => syncSelectedSlideHeight(api));
								}
							}}
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
		<FileDisplay username={username} post={{ postType, postOwner } as ScrapeResponse} file={files[0]} />
	);
}

export function FileDisplay({
	username,
	file,
	post: { postType, postOwner },
	onMediaLoad,
}: {
	username: string;
	file: string;
	post: ScrapeResponse;
	onMediaLoad?: () => void;
}) {
	const url = `/api/storage/${username}/${postTypeString(postType)}/${postOwner}/${file}`;
	if (/\.(jpg)|(jpeg)|(webp)|(heic)$/.test(file)) {
		return <img src={url} onLoad={onMediaLoad} className="h-auto w-full" />;
	} else if (/\.(mp4)|(webm)$/.test(file)) {
		return <video src={url} loop controls muted onLoadedMetadata={onMediaLoad} className="h-auto w-full" />;
	} else {
		return <a href={url}>{url}</a>;
	}
}
