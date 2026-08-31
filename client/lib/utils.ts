import { timestampDate, type Timestamp } from "@bufbuild/protobuf/wkt";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

import { PostType } from "@/buf/raker/v1/raker_pb";
import { toast, type ToastPosition } from "@/components/ui/toast";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export async function writeClipboard(text: string) {
	try {
		await navigator.clipboard.writeText(text);
		Toaster.success(`${text} was copied to your clipboard`);
	} catch (err) {
		Toaster.error(err as Error);
	}
}

export const defaultPostTypes = [
	PostType.Instagram,
	PostType.Highlight,
	PostType.Story,
	PostType.TikTok,
	PostType.Snapchat,
	PostType.VSCO,
];

export class Toaster {
	static success(message: string, position: ToastPosition = "top-center") {
		toast.add({
			title: message,
			type: "success",
			data: {
				position,
			},
		});
	}

	static error(err: Error, position: ToastPosition = "top-center") {
		toast.add({
			description: (err as Error).message,
			title: (err as Error).name,
			type: "error",
			data: {
				position,
			},
		});
	}
}

export function inPWA(): boolean {
	return (
		window.matchMedia("(display-mode: standalone)").matches ||
		window.matchMedia("(display-mode: fullscreen)").matches ||
		window.matchMedia("(display-mode: minimal-ui)").matches
	);
}

export function uniqueArraysEqualAsSets<T>(a: T[] = [], b: T[] = []) {
	if (a.length !== b.length) return false;
	const setA = new Set(a);
	return b.every((bi) => setA.has(bi));
}

export const dateFormatter = new Intl.DateTimeFormat(navigator.language, {
	weekday: "long",
	day: "2-digit",
	month: "2-digit",
	year: "numeric",
	hour: "2-digit",
	minute: "2-digit",
	second: "2-digit",
	hour12: true,
	timeZoneName: "longOffset",
});

export function timestampFormat(date: Timestamp): string {
	return dateFormatter.format(timestampDate(date));
}

enum FileTypePatterns {
	JPEG = "jpe?g",
	WebP = "webp",
	WebM = "webm",
	HEIC = "heic",
	MP4 = "mp4",
}

export const imageTypesRegexp = new RegExp(
	`\\.(${FileTypePatterns.JPEG})|(${FileTypePatterns.WebP})|(${FileTypePatterns.HEIC})\$`,
);
export const videoTypesRegexp = new RegExp(`\\.(${FileTypePatterns.MP4})|(${FileTypePatterns.WebM})\$`);
export const editFileTypesRegexp = new RegExp(
	`\\.(${FileTypePatterns.JPEG})|(${FileTypePatterns.WebP})|(${FileTypePatterns.MP4})\$`,
);
export const cropFileTypesRegexp = new RegExp(`\\.(${FileTypePatterns.JPEG})|(${FileTypePatterns.WebP})\$`);
