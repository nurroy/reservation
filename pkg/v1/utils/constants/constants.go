package constants

import "errors"

const (
	// ApplicationStatusKTPSave application status for success update Document Information
	ApplicationStatusKTPSave = "KTP_SAVED"

	// ApplicationStatusPerSave application status for success update Applicant Information
	ApplicationStatusPerSave = "PER_SAVED"

	// ApplicationStatusPaySave application status for success update Payroll Information
	ApplicationStatusPaySave = "PAY_SAVED"

	// ApplicationStatusConSave application status for success update Contact Information
	ApplicationStatusConSave = "CON_SAVED"

	// ApplicationStatusEmpSave application status for success update Employee Information
	ApplicationStatusEmpSave = "EMP_SAVED"

	// ApplicationStatusCreSub application status for Credit Score Submit success
	ApplicationStatusCreSub = "CR_SCORE_SUB"

	// ApplicationStatusCreRej application status for Credit Score Reject success
	ApplicationStatusCreRej = "CR_SCORE_REJ"

	// ApplicationStatusCreApr application status for Credit Score Approval success
	ApplicationStatusCreApr = "CR_SCORE_APR"

	// ApplicationStatusPayRej application status if the request payroll rejected
	ApplicationStatusPayRej = "PAYROLL_REJ"

	// ApplicationStatusUsrRej application status if the request rejected by user
	ApplicationStatusUsrRej = "USR_REJECT"

	// ApplicationStatusLoanCreated application status if the loan has been created
	ApplicationStatusLoanCreated = "LOAN_CREATED"

	// ApplicationStatusLoanPaid application status if the loan has been paid
	ApplicationStatusLoanPaid = "LOAN_PAID"

	// ApplicationStatusAppExp application status if the loan has expired
	ApplicationStatusAppExp = "APP_EXPIRED"

	// OCRTEXT for text detection type in OCR service
	OCRTEXT = "TEXT_DETECTION"

	// OCRVerFail to identify if OCR verification failed
	OCRVerFail = "OCR_VER_FAIL"

	// KTPReject to identify when the requested document rejected
	KTPReject = "KTP_REJECT"

	// ActLinkConf constanst application status ...
	ActLinkConf = "ACT_LINK_CONF"

	// DisbAccConf constanst application status disbursment confirmed
	DisbAccConf = "DISB_ACC_CONF"

	// EKYCompleteStage constanst application status E-KYC Completed
	EKYCompleteStage = "EKYC_COM"

	// PrivyCompleteStage constants for privy complete stage (DIG_SIGN_COM)
	PrivyCompleteStage = "DIG_SIGN_COM"

	// DocSigned constanst application status for document signing process complete
	DocSigned = "DOCUMENT_SIGNED"

	// DisbInProcess is the constants flag for starting the disbursement process
	DisbInProcess = "DISB_IN_PROCESS"

	// ISOTimeLayout ISO standard time layout without timezone (with timezone use time.RFC3339 instead)
	ISOTimeLayout = "2006-01-02T15:04:05"

	// XBRITimeLayout standard time layout with timezone
	XBRITimeLayout = "2006-01-02T15:04:05.000Z"

	// NoRecordsFound error message for no records found
	NoRecordsFound = "No records fetched"

	// UserIDIsInvalid error message for no records found
	UserIDIsInvalid = "User ID is invalid"

	// AccountIDIsInvalid error message for no records found
	AccountIDIsInvalid = "Account ID is invalid"

	// UserIDIsEmpty error message for no records found
	UserIDIsEmpty = "User ID is empty"

	// NoHostFound error message for not found host
	NoHostFound = "The host does not exist."

	// LoanStatusActive flag for active loan exist in finacle
	LoanStatusActive = "Active"

	// UnableToProcessRequestErrMessage error message for unable to process request
	UnableToProcessRequestErrMessage = "We are unable to process your request. Try after sometime."

	// UnexpectedErrorDuringRetrieval ...
	UnexpectedErrorDuringRetrieval = "An unexpected exception occurred during loan list retrieval"

	// UnexpectedErrorBalance balance error
	UnexpectedErrorBalance = "Error in Fetching in Balance details."

	// PinangProductName project product name
	PinangProductName = "PINANG"

	// FinacleXMLSchema ...
	FinacleXMLSchema = "http://www.finacle.com/fixml"

	// FinacleXMLSchemaInstance ...
	FinacleXMLSchemaInstance = "http://www.w3.org/2001/XMLSchema-instance"

	// LoanApplicationMissingOrInvalid error message for loan invalid or null
	LoanApplicationMissingOrInvalid = "Loan Application ID is NULL or Invalid"

	// PrivyRecordNotFoundErrMessage error message for invalid privy
	PrivyRecordNotFoundErrMessage = "Privy record not found"

	//InternalErrorASLIRIMessage error message for internal error ASLIRI
	InternalErrorASLIRIMessage = "Internal service error, try in 15 minutes."

	//ASTG code type constants
	ASTG = "ASTG"

	//FinacleURLPath url path for finacle services
	FinacleURLPath = "FISERVLET/fihttp"

	//CompareDone status after photo compare
	CompareDone = "COMP_DONE"

	// LoanTypeCode code type constants LOAN TYPE
	LoanTypeCode = "LNTP"

	// CommonCodePinang code type constants for cocd
	CommonCodePinang = "Pinang"

	// CommonCodeMaritalStatus code type constants for cocd
	CommonCodeMaritalStatus = "MSTH"

	// CommonCodeLastEducation code type constants for cocd last education
	CommonCodeLastEducation = "EDUH"

	// CommonCodeHomeOwnership code type constants for cocd home ownership
	CommonCodeHomeOwnership = "RESH"

	// CommonCodeHomeOwnership code type constants for cocd segmentation class
	CommonCodeSegmentationType = "SGMC"

	// DocSignedCodeDesc code desc constants for astg DOCUMENT_SIGNED
	DocSignedCodeDesc = "Document Signed"

	// EkycCompletedCodeDesc code desc constants for astg EKYC_COM
	EkycCompletedCodeDesc = "EKYC Completed"

	// PrivyCompleteStageCodeDesc code desc constants for astg DIG_SIGN_COM
	PrivyCompleteStageCodeDesc = "Digital Signature Created"

	// LoanCreatedCodeDesc code desc constants for astg LOAN_CREATED
	LoanCreatedCodeDesc = "Loan Created"
)

// Storing all the codes constant
const (
	// CodeLoanApplicationFetchNoRecords error code for record not found in the database
	CodeLoanApplicationFetchNoRecords = "200606"

	// CodeUserIDEmpty error code for user id empty
	CodeUserIDEmpty = "211081"

	// CodeUserIDInvalid error code for user id invalid / not found
	CodeUserIDInvalid = "211082"

	// CodeAccountIDInvalid error code for user id invalid / not found
	CodeAccountIDInvalid = "211083"

	// CodeHostNotFound error code for host invalid / not found
	CodeHostNotFound = "14084"

	// CodeHostNotFound error code for host invalid / not found
	CodeUnableProcessRequest = "100126"

	// CodeHostNotFound error code for host invalid / not found
	CodeBalanceProcessFailed = "211304"

	// SuccessCode is the constant code for api success
	SuccessCode = "0000"

	// CodeLoanAppIDMandatory
	CodeLoanAppIDMandatory = "211095"

	// CodeASLIRI02 asliri03 error code
	CodeASLIRI02 = "211102"

	// CodeASLIRI03 asliri03 error code
	CodeASLIRI03 = "211103"

	// CodeASLIRI04 asliri03 error code
	CodeASLIRI04 = "211104"

	// CodeASLIRI05 asliri03 error code
	CodeASLIRI05 = "211105"

	// CodeASLIRI06 asliri03 error code
	CodeASLIRI06 = "211106"
)

// Storing all the cm code desc constant
const (
	// CMDescInvalidInputValue cm code for invalid input message
	CMDescInvalidInputValue = "There are invalid input value"
)

// Constants Property
const (
	// PRPMAppId constant for Property APP ID
	PRPMAppID = "BWY"

	// CreditScoreTimeIntervalProperty property name for credit score time interval
	CreditScoreTimeIntervalProperty = "CREDIT_SCR_CAL_TIME_INTRVL"

	// CooldownPeriodMissingSimpananProperty property name for cooldown period missing simpanan
	CooldownPeriodMissingSimpananProperty = "COOLDOWN_PERIOD_FOR_MISSING_SIMPANAN"

	// CreditOfferValidityPeriodProperty property name for credit offer validity
	CreditOfferValidityPeriodProperty = "CREDIT_OFFER_VALIDITY_PERIOD"

	// InstallationProductIDProperty property name for installation product id
	InstallationProductIDProperty = "INSTALLATION_PRODUCT_ID"

	// AllowedPeriodIncompleteProperty property name for Allowed period in days to complete the incomplete credit approval process
	AllowedPeriodIncompleteProperty = "ALLOWED_PERIOD_FOR_INCOMPLETE_APP"

	// CooldownPeriodPayrollRejProperty property name for Cooldown period in days for a system rejected payroll application
	CooldownPeriodPayrollRejProperty = "COOLDOWN_PERIOD_FOR_PAYROLL_REJECTION"

	// AllowedConsecutiveUserRejProperty property name for Allowed number of user rejection post which a temporary ban will be put on user for application
	AllowedConsecutiveUserRejProperty = "ALLOWED_CONSECUTIVE_USER_REJ"

	// BanPeriodForUserRejProperty property name for Cooldown period in days for a user rejected application
	BanPeriodForUserRejProperty = "BAN_PERIOD_FOR_USER_REJECTION"

	// BRILifeInsuranceCode property name for "BRI_LIFE_INSURANCE_CODE"
	BRILifeInsuranceCode = "BRI_LIFE_INSURANCE_CODE"

	// FEBASharedSystemPath property name for FEBA_SHARED_SYS_PATH
	FEBASharedSystemPath = "FEBA_SHARED_SYS_PATH"

	// IsEKYCEnabled property name for IS_EKYC_ENABLED
	IsEKYCEnabled = "IS_EKYC_ENABLED"

	// AsliriApprovalRate property name for ASLIRI_APPROVAL_RATE
	AsliriApprovalRate = "ASLIRI_APPROVAL_RATE"

	// IsAsliriEnabled property name for IS_ASLIRI_ENABLED
	IsAsliriEnabled = "IS_ASLIRI_ENABLED"

	// ASLIRIURLRequest is prpm value for retrieve the verify URL ASLIRI_SELFIE_VERIFY_URL
	ASLIRIURLRequest = "ASLIRI_SELFIE_VERIFY_URL"

	// ASLIRITokenRequest is prpm value for retrieve the token ASLIRI_HEADER_TOKEN
	ASLIRITokenRequest = "ASLIRI_HEADER_TOKEN"

	//DaysToDueDateInstallmentOne DAYS_TO_DUE_DATE_INST_ONE
	DaysToDueDateInstallmentOne = "DAYS_TO_DUE_DATE_INST_ONE"

	//CoreBODDate CORE_BOD_DATE
	CoreBODDate = "CORE_BOD_DATE"

	//DaysToDueDateInstallmentTwoOrMore DAYS_TO_DUE_DATE_INST_TWO_OR_MORE
	DaysToDueDateInstallmentTwoOrMore = "DAYS_TO_DUE_DATE_INST_TWO_OR_MORE"

	//PrivyMerchantKey PRIVY_MERCHANT_KEY
	PrivyMerchantKey = "PRIVY_MERCHANT_KEY"

	//PrivyAuthToken PRIVY_AUTHORIZATION_KEY
	PrivyAuthToken = "PRIVY_AUTHORIZATION_KEY"

	//PersonalNumberKey PERSONAL_NUMBER
	PersonalNumberKey = "PERSONAL_NUMBER"
)

var (
	// ErrNoRecords error default if there is no records in the database
	ErrNoRecords = errors.New(NoRecordsFound)

	// ErrUserIDInvalid error default if there is invalid user ID provided
	ErrUserIDInvalid = errors.New(UserIDIsInvalid)

	// ErrUnexpectedDuringRetrieval error default if there is failure in finacle
	ErrUnexpectedDuringRetrieval = errors.New(UnexpectedErrorDuringRetrieval)

	// ErrProcessLoanCreation error default if there is failure in finacle create loan
	ErrProcessLoanCreation = errors.New("error while creating loan")

	// ErrProcessCustomerCreation error default if there is failure in finacle create customer
	ErrProcessCustomerCreation = errors.New("error while creating customer")

	// ErrProcessPrivyCreation error default if there is failure create privy registration
	ErrProcessPrivyCreation = errors.New("error while creating privy registration")

	// ErrRuntime error runtime
	ErrRuntime = errors.New("Runtime error has occurred")
)
