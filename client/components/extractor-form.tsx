import type { DescMethodUnary, MessageInitShape, MessageShape } from "@bufbuild/protobuf";
import { useMutation } from "@connectrpc/connect-query";
import type { AnyFieldApi } from "@tanstack/form-core";
import { useForm } from "@tanstack/react-form";
import type { StandardSchemaV1 } from "@tanstack/react-form";
import type { UseNavigateResult } from "@tanstack/router-core";
import { SprayCanIcon } from "lucide-react";
import { useEffect, useRef, useState, type Dispatch, type ReactNode, type SetStateAction } from "react";

import type { ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { Result } from "@/components/result";
import { Button } from "@/components/ui/button";
import { CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress";
import { Switch } from "@/components/ui/switch";
import { useUser } from "@/hooks/user-provider";
import { Toaster } from "@/lib/utils";

type ExtractorSearchValues = Record<string, string | boolean>;

type ExtractorSprayConfig<TMutation extends DescMethodUnary> = {
	enabled: boolean;
	isSuccessResult?: (result: MessageShape<TMutation["output"]>) => boolean;
	delayMs?: number;
};

type ExtractorFormConfig<
	TFrom extends string,
	TSearch extends ExtractorSearchValues,
	TMutation extends DescMethodUnary,
> = {
	navigate: UseNavigateResult<TFrom>;
	search: TSearch;
	validators: {
		onChange: StandardSchemaV1<TSearch, unknown>;
		onSubmit: StandardSchemaV1<TSearch, unknown>;
	};
	mutation: TMutation;
	autoSubmitWhen: (search: TSearch) => boolean;
	buildMutationArgs: (search: TSearch) => MessageInitShape<TMutation["input"]>;
	buildSearch: (search: TSearch, result: MessageShape<TMutation["output"]>) => TSearch;
	sprayConfig?: ExtractorSprayConfig<TMutation>;
};

export function useExtractorForm<
	TFrom extends string,
	TSearch extends ExtractorSearchValues,
	TMutation extends DescMethodUnary,
>({
	navigate,
	search,
	validators,
	mutation,
	autoSubmitWhen,
	buildMutationArgs,
	buildSearch,
	sprayConfig,
}: ExtractorFormConfig<TFrom, TSearch, TMutation>) {
	const { username } = useUser();
	const extractorMutation = useMutation<TMutation["input"], TMutation["output"]>(mutation);
	const [result, setResult] = useState<MessageShape<TMutation["output"]> | null>(null);
	const [sprayModeEnabled, setSprayModeEnabled] = useState(false);
	const [sprayTimes, setSprayTimes] = useState(3);

	const form = useForm({
		defaultValues: search,
		validators,
		onSubmit: async ({ value }) => {
			if (sprayModeEnabled && sprayConfig?.enabled) {
				// Spray mode: loop until success or max times reached
				let currentResult: MessageShape<TMutation["output"]> | null = null;
				const isSuccess = sprayConfig.isSuccessResult ?? (() => true);
				const delayMs = sprayConfig.delayMs ?? 500;

				for (let i = 0; i < sprayTimes; i++) {
					try {
						currentResult = await extractorMutation.mutateAsync(buildMutationArgs(value));

						if (isSuccess(currentResult)) {
							setResult(currentResult);
							await navigate({
								search: buildSearch(value, currentResult) as never,
								replace: true,
							});
							break;
						}

						// Add a small delay between attempts
						if (i < sprayTimes - 1) {
							await new Promise((resolve) => setTimeout(resolve, delayMs));
						}
					} catch (err) {
						Toaster.error(err as Error);
					}
				}
			} else {
				try {
					// Normal mode: single execution
					const result = await extractorMutation.mutateAsync(buildMutationArgs(value));
					setResult(result);
					await navigate({ search: buildSearch(value, result) as never, replace: true });
				} catch (err) {
					Toaster.error(err as Error);
				}
			}
		},
	});

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

	const initialSubmit = useRef(true);
	useEffect(() => {
		if (!initialSubmit.current) return;
		initialSubmit.current = false;
		if (username === null) return;
		if (autoSubmitWhen(search)) {
			form.handleSubmit();
		}
	}, [autoSubmitWhen, form, search, username]);

	const searchKey = JSON.stringify(search);
	useEffect(() => {
		for (const [key, value] of Object.entries(search)) {
			form.setFieldValue(key as never, value as never);
		}
	}, [form, search, searchKey]);

	return {
		form,
		result,
		setResult,
		isPending: extractorMutation.isPending,
		sprayModeEnabled,
		setSprayModeEnabled,
		sprayTimes,
		setSprayTimes,
		sprayConfigEnabled: sprayConfig?.enabled ?? false,
	};
}

export function ExtractorTextField({
	field,
	label,
	placeholder,
}: {
	field: AnyFieldApi;
	label: string;
	placeholder: string;
}) {
	const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;

	return (
		<Field>
			<FieldLabel htmlFor={field.name}>{label}</FieldLabel>
			<Input
				id={field.name}
				name={field.name}
				value={field.state.value}
				onBlur={field.handleBlur}
				aria-invalid={isInvalid}
				autoCapitalize="off"
				autoComplete="off"
				autoCorrect="off"
				onChange={(e) => field.handleChange(e.target.value)}
				placeholder={placeholder}
			/>
			{isInvalid && <FieldError errors={field.state.meta.errors} />}
		</Field>
	);
}

export function ExtractorSwitchField({ field, label }: { field: AnyFieldApi; label: string }) {
	return (
		<Field orientation="horizontal" className="w-fit">
			<FieldLabel htmlFor={field.name}>{label}</FieldLabel>
			<Switch
				id={field.name}
				name={field.name}
				checked={field.state.value}
				onCheckedChange={field.handleChange}
			/>
		</Field>
	);
}

export function ExtractorSprayControls({
	enabled,
	onEnabledChange,
	times,
	onTimesChange,
}: {
	enabled: boolean;
	onEnabledChange: (enabled: boolean) => void;
	times: number;
	onTimesChange: (times: number) => void;
}) {
	return (
		<div className="flex flex-row items-center gap-2 sm:gap-4">
			<Field orientation="horizontal" className="w-fit">
				<Switch id="spray-mode" checked={enabled} onCheckedChange={onEnabledChange} />
				<SprayCanIcon />
				{/* <FieldLabel htmlFor="spray-mode">Spray Mode</FieldLabel> */}
			</Field>
			<Field className="w-fit">
				{/* <FieldLabel htmlFor="spray-times">Spray Times</FieldLabel> */}
				<Input
					id="spray-times"
					type="number"
					min="1"
					max="100"
					value={times}
					onChange={(e) => onTimesChange(Math.max(1, parseInt(e.target.value) || 1))}
					disabled={!enabled}
					placeholder="Number of attempts"
					className="w-20"
				/>
			</Field>
		</div>
	);
}

export function ExtractorFormShell({
	form,
	isPending,
	result,
	setResult,
	children,
	sprayConfigEnabled = false,
	sprayModeEnabled = false,
	sprayTimes = 3,
	onSprayModeChange,
	onSprayTimesChange,
}: {
	form: { handleSubmit: () => void };
	isPending: boolean;
	result: ScrapeResponse | null;
	setResult: Dispatch<SetStateAction<ScrapeResponse | null>>;
	children: ReactNode;
	sprayConfigEnabled?: boolean;
	sprayModeEnabled?: boolean;
	sprayTimes?: number;
	onSprayModeChange?: (enabled: boolean) => void;
	onSprayTimesChange?: (times: number) => void;
}) {
	return (
		<form
			onSubmit={(e) => {
				e.preventDefault();
				form.handleSubmit();
			}}
		>
			<CardContent>
				<FieldGroup>
					{children}
					{sprayConfigEnabled && (
						<ExtractorSprayControls
							enabled={sprayModeEnabled}
							onEnabledChange={onSprayModeChange || (() => {})}
							times={sprayTimes}
							onTimesChange={onSprayTimesChange || (() => {})}
						/>
					)}
					<Field orientation="horizontal">
						<Button type="submit" disabled={isPending} className="mb-3 w-full sm:w-auto">
							Submit
						</Button>
					</Field>
					{isPending && (
						<Field>
							<Progress value={null} className="pb-2" />
						</Field>
					)}
				</FieldGroup>
			</CardContent>
			{result && (
				<CardFooter>
					<Result result={result} setResult={setResult} />
				</CardFooter>
			)}
		</form>
	);
}
