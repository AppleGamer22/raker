import { Trash2Icon } from "lucide-react";
import { useState, useCallback, useRef } from "react";

import {
	AlertDialog,
	AlertDialogContent,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogMedia,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";

interface ConfirmationRequest {
	title: string;
	description: string;
	onConfirm: () => void;
	onCancel: () => void;
	cancelText?: string;
	confirmText?: string;
	isDestructive?: boolean;
}

export function useConfirmationDialog() {
	const [isOpen, setIsOpen] = useState(false);
	const requestRef = useRef<ConfirmationRequest | null>(null);

	const confirm = useCallback((params: Omit<ConfirmationRequest, "onConfirm" | "onCancel">): Promise<boolean> => {
		return new Promise((resolve) => {
			requestRef.current = {
				...params,
				onConfirm: () => {
					setIsOpen(false);
					resolve(true);
				},
				onCancel: () => {
					setIsOpen(false);
					resolve(false);
				},
			};
			setIsOpen(true);
		});
	}, []);

	const handler = requestRef.current;

	const DialogComponent = () => (
		<AlertDialog open={isOpen} onOpenChange={(open) => !open && handler?.onCancel?.()}>
			<AlertDialogContent size="sm">
				<AlertDialogHeader>
					<AlertDialogMedia className="bg-destructive/10 text-destructive dark:bg-destructive/20 dark:text-destructive">
						<Trash2Icon />
					</AlertDialogMedia>
					<AlertDialogTitle>{handler?.title}</AlertDialogTitle>
					{handler?.description && <AlertDialogDescription>{handler.description}</AlertDialogDescription>}
				</AlertDialogHeader>
				<AlertDialogFooter>
					<Button variant="outline" onClick={handler?.onCancel}>
						{handler?.cancelText ?? "Cancel"}
					</Button>
					<Button variant={handler?.isDestructive ? "destructive" : "default"} onClick={handler?.onConfirm}>
						{handler?.confirmText ?? "Confirm"}
					</Button>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	);

	return {
		confirm,
		DialogComponent,
	};
}
