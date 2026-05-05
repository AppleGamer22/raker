import { clsx, type ClassValue } from "clsx";
import { toast } from "sonner";
import { twMerge } from "tailwind-merge";

import { PostType } from "@/buf/raker/v1/raker_pb";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export async function writeClipboard(text: string) {
	try {
		await navigator.clipboard.writeText(text);
		toast.success(`${text} was copied to your clipboard`, {
			position: "top-center",
		});
	} catch (err) {
		toast.error((err as Error).message, {
			position: "top-center",
		});
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

export function inPWA(): boolean {
	return (
		window.matchMedia("(display-mode: standalone)").matches ||
		window.matchMedia("(display-mode: fullscreen)").matches ||
		window.matchMedia("(display-mode: minimal-ui)").matches
	);
}
