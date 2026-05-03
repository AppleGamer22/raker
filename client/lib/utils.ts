import { clsx, type ClassValue } from "clsx";
import { toast } from "sonner";
import { twMerge } from "tailwind-merge";

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
