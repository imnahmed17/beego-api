$('.t-datepicker').tDatePicker({

    // auto close after selection
    autoClose        : true,

    // animation speed in milliseconds
    durationArrowTop : 200,

    // the number of calendars
    numCalendar    : 2,

    // localization
    titleCheckIn   : 'Check In',
    titleCheckOut  : 'Check Out',
    titleToday     : 'Today',
    titleDateRange : 'night',
    titleDateRanges: 'nights',
    titleDays      : [ 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa', 'Su' ],
    titleMonths    : ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'Septemper', 'October', 'November', "December"],

    // the max length of the title
    titleMonthsLimitShow : 3,

    // replace moth names
    replaceTitleMonths : null,

    // e.g. 'dd-mm-yy'
    showDateTheme   : null,

    // icon options
    iconArrowTop : true,
    iconDate     : '<i class="fi fi-rr-calendar-lines"></i>',
    arrowPrev    : '&#x276E;',
    arrowNext    : '&#x276F;',

    // shows today title
    toDayShowTitle       : true, 

    // shows date range title
    dateRangesShowTitle  : true,

    // highlights today
    toDayHighlighted     : false,

    // highlights next day
    nextDayHighlighted   : false,

    // an array of days
    daysOfWeekHighlighted: [0,6],

    // custom date format
    formatDate      : 'yyyy-mm-dd',

    dateCheckIn  : new Date(),
    dateCheckOut : new Date(new Date().getTime() + 24 * 60 * 60 * 1000),
    startDate    : null,
    endDate      : null,

    // limits the number of months
    limitPrevMonth : 0,
    limitNextMonth : 11,

    // limits the number of days
    limitDateRanges    : 31,

    // true -> full days || false - 1 day
    showFullDateRanges : false, 

    // DATA HOLIDAYS
    // Data holidays
    fnDataEvent   : null
});