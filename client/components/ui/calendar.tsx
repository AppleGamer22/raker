import { cn } from "cn";
import { addMonths, endOfDay, format, isAfter, isSameMonth, isValid, parse } from "date-fns";
import { ChevronLeftIcon, ChevronRightIcon, ChevronDownIcon, XIcon } from "lucide-react";
import { CalendarIcon } from "lucide-react";
import * as React from "react";
import { DayPicker, getDefaultClassNames, type DayButton, type Locale } from "react-day-picker";
import { type DateRange } from "react-day-picker";

import { Button, buttonVariants } from "@/components/ui/button";
import { Field, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { dateOnlyFormatter, formatDateOnly } from "@/lib/utils";

export function Calendar({
	className,
	classNames,
	showOutsideDays = true,
	captionLayout = "label",
	buttonVariant = "ghost",
	locale,
	formatters,
	components,
	...props
}: React.ComponentProps<typeof DayPicker> & {
	buttonVariant?: React.ComponentProps<typeof Button>["variant"];
}) {
	const defaultClassNames = getDefaultClassNames();

	return (
		<DayPicker
			showOutsideDays={showOutsideDays}
			className={cn(
				"group/calendar bg-background p-2 [--cell-radius:var(--radius-md)] [--cell-size:--spacing(7)] in-data-[slot=card-content]:bg-transparent in-data-[slot=popover-content]:bg-transparent",
				String.raw`rtl:**:[.rdp-button\_next>svg]:rotate-180`,
				String.raw`rtl:**:[.rdp-button\_previous>svg]:rotate-180`,
				className,
			)}
			captionLayout={captionLayout}
			locale={locale}
			formatters={{
				formatMonthDropdown: (date) => date.toLocaleString(locale?.code, { month: "short" }),
				...formatters,
			}}
			classNames={{
				root: cn("w-fit", defaultClassNames.root),
				months: cn("relative flex flex-col gap-4 sm:flex-row", defaultClassNames.months),
				month: cn("flex w-full flex-col gap-4", defaultClassNames.month),
				nav: cn(
					"absolute inset-x-0 top-0 flex w-full items-center justify-between gap-1",
					defaultClassNames.nav,
				),
				button_previous: cn(
					buttonVariants({ variant: buttonVariant }),
					"size-(--cell-size) p-0 select-none aria-disabled:opacity-50",
					defaultClassNames.button_previous,
				),
				button_next: cn(
					buttonVariants({ variant: buttonVariant }),
					"size-(--cell-size) p-0 select-none aria-disabled:opacity-50",
					defaultClassNames.button_next,
				),
				month_caption: cn(
					"flex h-(--cell-size) w-full items-center justify-center px-(--cell-size)",
					defaultClassNames.month_caption,
				),
				dropdowns: cn(
					"flex h-(--cell-size) w-full items-center justify-center gap-1.5 text-sm font-medium",
					defaultClassNames.dropdowns,
				),
				dropdown_root: cn("relative rounded-(--cell-radius)", defaultClassNames.dropdown_root),
				dropdown: cn("absolute inset-0 bg-popover opacity-0", defaultClassNames.dropdown),
				caption_label: cn(
					"font-medium select-none",
					captionLayout === "label"
						? "text-sm"
						: "flex items-center gap-1 rounded-(--cell-radius) text-sm [&>svg]:size-3.5 [&>svg]:text-muted-foreground",
					defaultClassNames.caption_label,
				),
				month_grid: cn("w-full border-collapse", defaultClassNames.month_grid),
				weekdays: cn("flex", defaultClassNames.weekdays),
				weekday: cn(
					"flex-1 rounded-(--cell-radius) text-[0.8rem] font-normal text-muted-foreground select-none",
					defaultClassNames.weekday,
				),
				week: cn("mt-2 flex w-full", defaultClassNames.week),
				week_number_header: cn("w-(--cell-size) select-none", defaultClassNames.week_number_header),
				week_number: cn("text-[0.8rem] text-muted-foreground select-none", defaultClassNames.week_number),
				day: cn(
					"group/day relative aspect-square h-full w-full rounded-(--cell-radius) p-0 text-center select-none [&:last-child[data-selected=true]_button]:rounded-r-(--cell-radius)",
					props.showWeekNumber
						? "[&:nth-child(2)[data-selected=true]_button]:rounded-l-(--cell-radius)"
						: "[&:first-child[data-selected=true]_button]:rounded-l-(--cell-radius)",
					defaultClassNames.day,
				),
				range_start: cn(
					"relative isolate z-0 rounded-l-(--cell-radius) bg-muted after:absolute after:inset-y-0 after:right-0 after:w-4 after:bg-muted",
					defaultClassNames.range_start,
				),
				range_middle: cn("rounded-none", defaultClassNames.range_middle),
				range_end: cn(
					"relative isolate z-0 rounded-r-(--cell-radius) bg-muted after:absolute after:inset-y-0 after:left-0 after:w-4 after:bg-muted",
					defaultClassNames.range_end,
				),
				today: cn(
					"rounded-(--cell-radius) bg-muted text-foreground data-[selected=true]:rounded-none",
					defaultClassNames.today,
				),
				outside: cn("text-muted-foreground aria-selected:text-muted-foreground", defaultClassNames.outside),
				disabled: cn("text-muted-foreground opacity-50", defaultClassNames.disabled),
				hidden: cn("invisible", defaultClassNames.hidden),
				...classNames,
			}}
			components={{
				Root: ({ className, rootRef, ...props }) => {
					return <div data-slot="calendar" ref={rootRef} className={cn(className)} {...props} />;
				},
				Chevron: ({ className, orientation, ...props }) => {
					if (orientation === "left") {
						return <ChevronLeftIcon className={cn("size-4", className)} {...props} />;
					}

					if (orientation === "right") {
						return <ChevronRightIcon className={cn("size-4", className)} {...props} />;
					}

					return <ChevronDownIcon className={cn("size-4", className)} {...props} />;
				},
				DayButton: ({ ...props }) => <CalendarDayButton locale={locale} {...props} />,
				WeekNumber: ({ children, ...props }) => {
					return (
						<td {...props}>
							<div className="flex size-(--cell-size) items-center justify-center text-center">
								{children}
							</div>
						</td>
					);
				},
				...components,
			}}
			{...props}
		/>
	);
}

export function CalendarDayButton({
	className,
	day,
	modifiers,
	locale,
	...props
}: React.ComponentProps<typeof DayButton> & { locale?: Partial<Locale> }) {
	const defaultClassNames = getDefaultClassNames();

	const ref = React.useRef<HTMLButtonElement>(null);
	React.useEffect(() => {
		if (modifiers.focused) ref.current?.focus();
	}, [modifiers.focused]);

	return (
		<Button
			variant="ghost"
			size="icon"
			data-day={day.date.toLocaleDateString(locale?.code)}
			data-selected-single={
				modifiers.selected && !modifiers.range_start && !modifiers.range_end && !modifiers.range_middle
			}
			data-range-start={modifiers.range_start}
			data-range-end={modifiers.range_end}
			data-range-middle={modifiers.range_middle}
			className={cn(
				"relative isolate z-10 flex aspect-square size-auto w-full min-w-(--cell-size) flex-col gap-1 border-0 leading-none font-normal group-data-[focused=true]/day:relative group-data-[focused=true]/day:z-10 group-data-[focused=true]/day:border-ring group-data-[focused=true]/day:ring-[3px] group-data-[focused=true]/day:ring-ring/50 data-[range-end=true]:rounded-(--cell-radius) data-[range-end=true]:rounded-r-(--cell-radius) data-[range-end=true]:bg-primary data-[range-end=true]:text-primary-foreground data-[range-middle=true]:rounded-none data-[range-middle=true]:bg-muted data-[range-middle=true]:text-foreground data-[range-start=true]:rounded-(--cell-radius) data-[range-start=true]:rounded-l-(--cell-radius) data-[range-start=true]:bg-primary data-[range-start=true]:text-primary-foreground data-[selected-single=true]:bg-primary data-[selected-single=true]:text-primary-foreground dark:hover:text-foreground [&>span]:text-xs [&>span]:opacity-70",
				defaultClassNames.day,
				className,
			)}
			{...props}
		/>
	);
}

function parseDateInput(text: string): Date | undefined {
	const trimmed = text.trim();
	if (!trimmed) return undefined;

	// Check whether day precedes month in user's locale (e.g. en-AU: dd/mm/yyyy)
	const parts = dateOnlyFormatter.formatToParts(new Date(2026, 0, 1));
	const isDayFirst = parts.findIndex((p) => p.type === "day") < parts.findIndex((p) => p.type === "month");

	const localeFormats = isDayFirst
		? [
				"dd/MM/yyyy",
				"d/M/yyyy",
				"dd/MM/yy",
				"d/M/yy",
				"dd-MM-yyyy",
				"d-M-yyyy",
				"dd.MM.yyyy",
				"d.M.yyyy",
				"MM/dd/yyyy",
				"M/d/yyyy",
			]
		: ["MM/dd/yyyy", "M/d/yyyy", "MM/dd/yy", "M/d/yy", "MM-dd-yyyy", "M-d-yyyy", "dd/MM/yyyy", "d/M/yyyy"];

	const otherFormats = [
		"yyyy-MM-dd",
		"yyyy-M-d",
		"yyyy/MM/dd",
		"yyyy/M/d",
		"LLL dd, y",
		"LLL d, y",
		"LLLL dd, y",
		"LLLL d, y",
		"dd LLL yyyy",
		"d LLL yyyy",
		"dd LLLL yyyy",
		"d LLLL yyyy",
	];

	for (const fmt of [...localeFormats, ...otherFormats]) {
		const parsed = parse(trimmed, fmt, new Date());
		if (isValid(parsed) && parsed.getFullYear() > 1900 && parsed.getFullYear() < 2100) {
			return parsed;
		}
	}

	const fallback = new Date(trimmed);
	if (
		isValid(fallback) &&
		!isNaN(fallback.getTime()) &&
		fallback.getFullYear() > 1900 &&
		fallback.getFullYear() < 2100
	) {
		return fallback;
	}

	return undefined;
}

function useMediaQuery(query: string): boolean {
	const subscribe = React.useCallback(
		(callback: () => void) => {
			const media = window.matchMedia(query);
			media.addEventListener("change", callback);
			return () => media.removeEventListener("change", callback);
		},
		[query],
	);

	const getSnapshot = () => window.matchMedia(query).matches;
	const getServerSnapshot = () => false;

	return React.useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);
}

export function DatePickerRange({
	value,
	onChange,
	label = "Date Picker Range",
	id = "date-picker-range",
	className = "w-fit",
}: {
	value?: DateRange;
	onChange?: (dates: DateRange | undefined) => void;
	label?: string;
	id?: string;
	className?: string;
}) {
	const isSmallScreen = useMediaQuery("(max-width: 640px)");
	const [month, setMonth] = React.useState<Date | undefined>(value?.from ?? new Date());

	const fromTime = value?.from?.getTime();
	const toTime = value?.to?.getTime();
	const [prevFromTime, setPrevFromTime] = React.useState(fromTime);
	const [prevToTime, setPrevToTime] = React.useState(toTime);

	const [fromText, setFromText] = React.useState(value?.from ? formatDateOnly(value.from) : "");
	const [toText, setToText] = React.useState(value?.to ? formatDateOnly(value.to) : "");

	if (fromTime !== prevFromTime) {
		setPrevFromTime(fromTime);
		setFromText(value?.from ? formatDateOnly(value.from) : "");
	}

	if (toTime !== prevToTime) {
		setPrevToTime(toTime);
		setToText(value?.to ? formatDateOnly(value.to) : "");
	}

	// oxlint-disable-next-line react/purity
	const today = new Date();
	const currentMonth = month ?? today;
	const isCurrentMonthInView = isSmallScreen
		? isSameMonth(today, currentMonth)
		: isSameMonth(today, currentMonth) || isSameMonth(today, addMonths(currentMonth, 1));

	const commitFrom = () => {
		if (!fromText.trim()) {
			if (value?.from) {
				onChange?.(value.to ? { from: undefined, to: value.to } : undefined);
			}
			return;
		}
		const parsed = parseDateInput(fromText);
		if (parsed && !isAfter(parsed, endOfDay(today))) {
			if (value?.to && isAfter(parsed, value.to)) {
				setFromText(formatDateOnly(value.to));
				setToText(formatDateOnly(parsed));
				setMonth(value.to);
				onChange?.({ from: value.to, to: parsed });
			} else {
				setFromText(formatDateOnly(parsed));
				setMonth(parsed);
				onChange?.({ from: parsed, to: value?.to });
			}
		} else {
			setFromText(value?.from ? formatDateOnly(value.from) : "");
		}
	};

	const commitTo = () => {
		if (!toText.trim()) {
			if (value?.to) {
				onChange?.(value.from ? { from: value.from, to: undefined } : undefined);
			}
			return;
		}
		const parsed = parseDateInput(toText);
		if (parsed && !isAfter(parsed, endOfDay(today))) {
			if (value?.from && isAfter(value.from, parsed)) {
				setFromText(formatDateOnly(parsed));
				setToText(formatDateOnly(value.from));
				setMonth(parsed);
				onChange?.({ from: parsed, to: value.from });
			} else {
				setToText(formatDateOnly(parsed));
				setMonth(parsed);
				onChange?.({ from: value?.from, to: parsed });
			}
		} else {
			setToText(value?.to ? formatDateOnly(value.to) : "");
		}
	};

	return (
		<Field className={className}>
			{label && <FieldLabel htmlFor={id}>{label}</FieldLabel>}
			<Popover
				onOpenChange={(open) => {
					if (open) {
						setMonth(value?.from ?? new Date());
						setFromText(value?.from ? formatDateOnly(value.from) : "");
						setToText(value?.to ? formatDateOnly(value.to) : "");
					}
				}}
			>
				<PopoverTrigger
					render={
						<Button
							variant="outline"
							id={id}
							className="max-w-full justify-start truncate px-2.5 font-normal"
						>
							<CalendarIcon data-icon="inline-start" />
							{value?.from ? (
								value.to ? (
									<span className="truncate">
										{format(value.from, "LLL dd, y")} - {format(value.to, "LLL dd, y")}
									</span>
								) : (
									<span className="truncate">{format(value.from, "LLL dd, y")}</span>
								)
							) : (
								<span>Pick a date</span>
							)}
						</Button>
					}
				/>
				<PopoverContent className="w-auto max-w-[calc(100vw-1.5rem)] overflow-x-auto p-0" align="start">
					<div className="flex flex-col items-stretch gap-2 border-b px-3 py-2.5 sm:flex-row sm:items-center sm:gap-4 sm:px-4">
						<div className="min-w-0 flex-1">
							<Input
								id={`${id}-from`}
								aria-label="Start date"
								placeholder="Start date"
								className="h-8 w-full min-w-0 text-xs"
								value={fromText}
								onChange={(e) => setFromText(e.target.value)}
								onBlur={commitFrom}
								onKeyDown={(e) => {
									if (e.key === "Enter") {
										e.preventDefault();
										commitFrom();
									}
								}}
							/>
						</div>
						<div className="min-w-0 flex-1">
							<Input
								id={`${id}-to`}
								aria-label="End date"
								placeholder="End date"
								className="h-8 w-full min-w-0 text-xs"
								value={toText}
								onChange={(e) => setToText(e.target.value)}
								onBlur={commitTo}
								onKeyDown={(e) => {
									if (e.key === "Enter") {
										e.preventDefault();
										commitTo();
									}
								}}
							/>
						</div>
					</div>
					<Calendar
						mode="range"
						month={month}
						onMonthChange={setMonth}
						selected={value}
						onSelect={onChange}
						numberOfMonths={isSmallScreen ? 1 : 2}
						className="w-full max-w-full"
						disabled={[
							{
								// oxlint-disable-next-line react/purity
								after: new Date(),
							},
						]}
					/>
					<div className="flex items-center justify-between border-t p-2">
						<Button
							variant="ghost"
							size="sm"
							type="button"
							disabled={!value?.from && !value?.to && !fromText && !toText}
							onClick={() => {
								onChange?.({ from: undefined, to: undefined });
								setFromText("");
								setToText("");
							}}
						>
							<XIcon />
							Clear
						</Button>
						<Button
							variant="outline"
							size="sm"
							type="button"
							disabled={isCurrentMonthInView}
							onClick={() => {
								if (!isCurrentMonthInView) {
									setMonth(new Date());
								}
							}}
						>
							Today's Month
						</Button>
					</div>
				</PopoverContent>
			</Popover>
		</Field>
	);
}
