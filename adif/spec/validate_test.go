// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

import "testing"

type validateTest struct {
	field Field
	value string
	want  Validity
}

var emptyCtx ValidationContext

func testValidator(t *testing.T, tc validateTest, ctx ValidationContext, funcname string) {
	t.Helper()
	v := TypeValidators[tc.field.Type.Name]
	if got := v(tc.value, tc.field, ctx); got.Validity != tc.want {
		if got.Validity == Valid {
			t.Errorf("%s(%q, %s, ctx) got Valid, want %s", funcname, tc.value, tc.field.Name, tc.want)
		} else {
			t.Errorf("%s(%q, %s, ctx) want %s got %s %s", funcname, tc.value, tc.field.Name, tc.want, got.Validity, got.Message)
		}
	}
}

func TestValidateBoolean(t *testing.T) {
	tests := []validateTest{
		{field: QsoRandomField, value: "Y", want: Valid},
		{field: SilentKeyField, value: "y", want: Valid},
		{field: ForceInitField, value: "N", want: Valid},
		{field: SwlField, value: "n", want: Valid},
		{field: QsoRandomField, value: "YES", want: InvalidError},
		{field: SilentKeyField, value: "true", want: InvalidError},
		{field: ForceInitField, value: "F", want: InvalidError},
		{field: SwlField, value: "false", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateBoolean")
	}
}

func TestValidateNumber(t *testing.T) {
	tests := []validateTest{
		{field: AgeField, value: "120", want: Valid},
		{field: AgeField, value: "-0", want: Valid},
		{field: AntElField, value: "90", want: Valid},
		{field: AntElField, value: "-90", want: Valid},
		{field: AIndexField, value: "123.0", want: Valid},
		{field: DistanceField, value: "9876.", want: Valid},
		{field: DistanceField, value: "1234567890", want: Valid},
		{field: FreqField, value: "146.520000001", want: Valid},
		{field: MaxBurstsField, value: "0", want: Valid},
		{field: MaxBurstsField, value: "00", want: Valid},
		{field: MyAltitudeField, value: "-1234.56789", want: Valid},
		{field: RxPwrField, value: ".7", want: Valid},
		{field: TxPwrField, value: "1499.999", want: Valid},
		{field: AgeField, value: "121", want: InvalidError},
		{field: AgeField, value: "-1", want: InvalidError},
		{field: AntElField, value: "--30", want: InvalidError},
		{field: AntElField, value: "99", want: InvalidError},
		{field: AntElField, value: "-91", want: InvalidError},
		{field: AntElField, value: "2π", want: InvalidError},
		{field: AIndexField, value: "-0.1", want: InvalidError},
		{field: AIndexField, value: "100-1", want: InvalidError},
		{field: AIndexField, value: "420", want: InvalidError},
		{field: DistanceField, value: "١٢٣", want: InvalidError},
		{field: DistanceField, value: "+9876", want: InvalidError},
		{field: DistanceField, value: "1 234", want: InvalidError},
		{field: DistanceField, value: "1,234", want: InvalidError},
		{field: DistanceField, value: "-0.00000001", want: InvalidError},
		{field: FreqField, value: "1.4652e2", want: InvalidError},
		{field: FreqField, value: "7.074-7", want: InvalidError},
		{field: MaxBurstsField, value: "1.2.3", want: InvalidError},
		{field: MaxBurstsField, value: "NaN", want: InvalidError},
		{field: MyAltitudeField, value: "〸", want: InvalidError},
		{field: MyAltitudeField, value: "⁷", want: InvalidError},
		{field: RxPwrField, value: "", want: InvalidError},
		{field: TxPwrField, value: ".", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateNumber")
	}
}

func TestValidateInteger(t *testing.T) {
	// currently all IntegerDataType fields have a minimum >= 0
	tests := []validateTest{
		{field: StxField, value: "0", want: Valid},
		{field: StxField, value: "1234567890", want: Valid},
		{field: NrBurstsField, value: "98765432123456789", want: Valid},
		{field: SfiField, value: "123", want: Valid},
		{field: KIndexField, value: "0", want: Valid},
		{field: KIndexField, value: "5", want: Valid},
		{field: KIndexField, value: "9", want: Valid},
		{field: StxField, value: "1,234", want: InvalidError},
		{field: StxField, value: "-1", want: InvalidError},
		{field: StxField, value: "7thirty", want: InvalidError},
		{field: StxField, value: "III", want: InvalidError},
		{field: NrBurstsField, value: "-1234", want: InvalidError},
		{field: NrBurstsField, value: "", want: InvalidError},
		{field: NrBurstsField, value: "twenty", want: InvalidError},
		{field: NrBurstsField, value: "௮", want: InvalidError},
		{field: NrBurstsField, value: "Ⅺ", want: InvalidError},
		{field: SfiField, value: "301", want: InvalidError},
		{field: SfiField, value: "7F", want: InvalidError},
		{field: SfiField, value: "0x20", want: InvalidError},
		{field: KIndexField, value: "10", want: InvalidError},
		{field: KIndexField, value: "-5", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateInteger")
	}
}

func TestValidatePositiveInteger(t *testing.T) {
	tests := []validateTest{
		{field: CqzField, value: "1", want: Valid},
		{field: CqzField, value: "40", want: Valid},
		{field: TenTenField, value: "1010101010", want: Valid},
		{field: FistsField, value: "1", want: Valid},
		{field: FistsField, value: "0987654321", want: Valid},
		{field: MyIotaIslandIdField, value: "1", want: Valid},
		{field: MyIotaIslandIdField, value: "666", want: Valid},
		{field: MyIotaIslandIdField, value: "99999999", want: Valid},
		{field: ItuzField, value: "1", want: Valid},
		{field: ItuzField, value: "90", want: Valid},
		{field: CqzField, value: "0", want: InvalidError},
		{field: CqzField, value: "42", want: InvalidError},
		{field: TenTenField, value: "0", want: InvalidError},
		{field: TenTenField, value: "-123", want: InvalidError},
		{field: FistsField, value: "five", want: InvalidError},
		{field: FistsField, value: "", want: InvalidError},
		{field: MyIotaIslandIdField, value: "0", want: InvalidError},
		{field: MyIotaIslandIdField, value: "100000000", want: InvalidError},
		{field: ItuzField, value: "111", want: InvalidError},
		{field: ItuzField, value: "99", want: InvalidError},
		{field: ItuzField, value: "８", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidatePositiveInteger")
	}
}

func TestValidateDate(t *testing.T) {
	tests := []validateTest{
		{field: QsoDateField, value: "19300101", want: Valid},
		{field: QsoDateOffField, value: "20200317", want: Valid},
		{field: QslrdateField, value: "19991231", want: Valid},
		{field: QslsdateField, value: "20000229", want: Valid},
		{field: QrzcomQsoUploadDateField, value: "21000101", want: Valid},
		{field: LotwQslrdateField, value: "23450607", want: Valid},
		{field: QsoDateField, value: "19000101", want: InvalidError},
		{field: QsoDateField, value: "19800100", want: InvalidError},
		{field: QsoDateOffField, value: "202012", want: InvalidError},
		{field: QsoDateOffField, value: "21000229", want: InvalidError},
		{field: QslrdateField, value: "1031", want: InvalidError},
		{field: QslsdateField, value: "2001-02-03", want: InvalidError},
		{field: QrzcomQsoUploadDateField, value: "01/02/2003", want: InvalidError},
		{field: LotwQslrdateField, value: "01022003", want: InvalidError},
		{field: LotwQslsdateField, value: "20220431", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateDate")
	}
}

func TestValidateTime(t *testing.T) {
	tests := []validateTest{
		{field: TimeOnField, value: "0000", want: Valid},
		{field: TimeOffField, value: "000000", want: Valid},
		{field: TimeOnField, value: "1234", want: Valid},
		{field: TimeOffField, value: "123456", want: Valid},
		{field: TimeOnField, value: "235959", want: Valid},
		{field: TimeOffField, value: "2000", want: Valid},
		{field: TimeOnField, value: "012345", want: Valid},
		{field: TimeOffField, value: "0159", want: Valid},
		{field: TimeOnField, value: "00:00", want: InvalidError},
		{field: TimeOffField, value: "00:00:00", want: InvalidError},
		{field: TimeOnField, value: "1234pm", want: InvalidError},
		{field: TimeOffField, value: "12.34.56", want: InvalidError},
		{field: TimeOnField, value: "735", want: InvalidError},
		{field: TimeOffField, value: "12345", want: InvalidError},
		{field: TimeOnField, value: "0630 AM", want: InvalidError},
		{field: TimeOffField, value: "noon", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateTime")
	}
}

func TestValidateString(t *testing.T) {
	tests := []validateTest{
		{field: ProgramidField, value: "Log & Operate", want: Valid},
		{field: CallField, value: "P5K2JI", want: Valid},
		{field: CommentField, value: `~!@#$%^&*()_+=-{}[]|\:;"'<>,.?/`, want: Valid},
		{field: EmailField, value: "flann.o-brien@example.com (Flann O`Brien)", want: Valid},
		{field: MyCityField, value: "", want: Valid},
		{field: ProgramidField, value: "🌲file", want: InvalidError},
		{field: CallField, value: "Я7ПД", want: InvalidError},
		{field: CommentField, value: `full width ．`, want: InvalidError},
		{field: EmailField, value: "thor@planet.earþ", want: InvalidError},
		{field: MyCityField, value: "\n", want: InvalidError},
		{field: CommentField, value: "Good\r\nchat", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateString")
	}
}

func TestValidateMultilineString(t *testing.T) {
	tests := []validateTest{
		{field: AddressField, value: "1600 Pennsylvania Ave\r\nWashington, DC 25000\r\n", want: Valid},
		{field: NotesField, value: "They were\non a boat", want: Valid},
		{field: QslmsgField, value: "\r", want: Valid},
		{field: RigField, value: "5 watts and a wire", want: Valid},
		{field: AddressField, value: "1600 Pennsylvania Ave\r\nWashington, ⏦ 25000\r\n", want: InvalidError},
		{field: NotesField, value: "Vertical\vtab", want: InvalidError},
		{field: QslmsgField, value: "🧶", want: InvalidError},
		{field: RigField, value: "Non‑breaking‑hyphen", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateMultilineString")
	}
}

func TestValidateIntlString(t *testing.T) {
	tests := []validateTest{
		{field: CommentIntlField, value: "१٢௩Ⅳ໖⁷၈🄊〸", want: Valid},
		{field: CountryIntlField, value: "🇭🇰", want: Valid},
		{field: MyAntennaIntlField, value: "null\0000character", want: Valid},
		{field: MyCityIntlField, value: "مكة المكرمة, Makkah the Honored", want: Valid},
		{field: MyCountryIntlField, value: "", want: Valid},
		{field: MyNameIntlField, value: "zero​width	space", want: Valid},
		{field: CommentIntlField, value: "new\nline", want: InvalidError},
		{field: CountryIntlField, value: "carriage\rreturn", want: InvalidError},
		{field: MyAntennaIntlField, value: "blank\r\n\r\nline", want: InvalidError},
		{field: MyCityIntlField, value: "line end\n", want: InvalidError},
		{field: MyCountryIntlField, value: "\r\n", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateIntlString")
	}
}

func TestValidateIntlMultilineString(t *testing.T) {
	tests := []validateTest{
		{field: AddressIntlField, value: "१٢௩Ⅳ໖⁷၈🄊〸", want: Valid},
		{field: NotesIntlField, value: "🇭🇰", want: Valid},
		{field: QslmsgIntlField, value: "null\0000character", want: Valid},
		{field: RigIntlField, value: "مكة المكرمة, Makkah the Honored", want: Valid},
		{field: AddressIntlField, value: "", want: Valid},
		{field: NotesIntlField, value: "zero​width	space", want: Valid},
		{field: QslmsgIntlField, value: "new\nline", want: Valid},
		{field: RigIntlField, value: "carriage\rreturn", want: Valid},
		{field: AddressIntlField, value: "blank\r\n\r\nline", want: Valid},
		{field: NotesIntlField, value: "line end\n", want: Valid},
		{field: QslmsgIntlField, value: "\r\n", want: Valid},
		{field: RigIntlField, value: "hello\tworld\r\n", want: Valid},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateIntlMultilineString")
	}
}

func TestValidateEnumeration(t *testing.T) {
	tests := []validateTest{
		{field: AntPathField, value: "G", want: Valid},
		{field: ArrlSectField, value: "wy", want: Valid},
		{field: BandField, value: "1.25m", want: Valid},
		{field: BandField, value: "13CM", want: Valid},
		{field: DxccField, value: "42", want: Valid},
		{field: QsoCompleteField, value: "?", want: Valid},
		{field: ContField, value: "na", want: Valid},
		{field: AntPathField, value: "X", want: InvalidError},
		{field: ArrlSectField, value: "TX", want: InvalidError},
		{field: ArrlSectField, value: "42", want: InvalidError},
		{field: BandField, value: "18m", want: InvalidError},
		{field: BandField, value: "130mm", want: InvalidError},
		{field: BandField, value: "G", want: InvalidError},
		{field: DxccField, value: "US", want: InvalidError},
		{field: QsoCompleteField, value: "!", want: InvalidError},
		{field: ContField, value: "Europe", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateEnumeration")
	}
}

func TestValidateEnumScope(t *testing.T) {
	tests := []struct {
		validateTest
		values map[string]string
	}{
		// empty string is valid
		{
			validateTest: validateTest{field: StateField, value: "", want: Valid},
			values:       map[string]string{},
		},
		// Wyoming, presumably
		{
			validateTest: validateTest{field: MyStateField, value: "WY", want: Valid},
			values:       map[string]string{"MY_STATE": "WY"},
		},
		// Yukon Territory, Canada
		{
			validateTest: validateTest{field: StateField, value: "YT", want: Valid},
			values:       map[string]string{"DXCC": "1", "STATE": "YT"},
		},
		// Sichuan, China
		{
			validateTest: validateTest{field: MyStateField, value: "SC", want: Valid},
			values:       map[string]string{"MY_DXCC": "318", "MY_STATE": "SC"},
		},
		// Nagano, Japan
		{
			validateTest: validateTest{field: StateField, value: "09", want: Valid},
			values:       map[string]string{"DXCC": "339", "STATE": "09"},
		},
		// Vranov nad Toplou, Slovak Republic
		{
			validateTest: validateTest{field: MyStateField, value: "vrt", want: Valid},
			values:       map[string]string{"MY_DXCC": "504", "MY_STATE": "vrt"},
		},
		// Permskaya oblast (new) or Permskaya Kraj (old) in Russia
		{
			validateTest: validateTest{field: StateField, value: "PM", want: Valid},
			values:       map[string]string{"DXCC": "15", "STATE": "PM"},
		},

		// not an empty string
		{
			validateTest: validateTest{field: StateField, value: "  ", want: InvalidError},
			values:       map[string]string{},
		},
		// Not a state abbreviation in any country
		{
			validateTest: validateTest{field: MyStateField, value: "XYZ", want: InvalidError},
			values:       map[string]string{"MY_STATE": "XYZ"},
		},
		// Not a valid DXCC code
		{
			validateTest: validateTest{field: MyStateField, value: "CA", want: InvalidWarning},
			values:       map[string]string{"MY_DXCC": "9876", "MY_STATE": "CA"},
		},
		// Yukon Territory, but wrong country
		{
			validateTest: validateTest{field: StateField, value: "YT", want: InvalidWarning},
			values:       map[string]string{"DXCC": "123", "STATE": "YT"},
		},
		// Sichuan, but not abbreviated
		{
			validateTest: validateTest{field: MyStateField, value: "Sichuan", want: InvalidError},
			values:       map[string]string{"MY_DXCC": "China", "MY_STATE": "Sichuan"},
		},
		// Nagano, but missing leading 0
		{
			validateTest: validateTest{field: StateField, value: "9", want: InvalidError},
			values:       map[string]string{"DXCC": "339", "STATE": "09"},
		},
		// Vranov nad Toplou, but Slovak Republic is spelled out
		{
			validateTest: validateTest{field: MyStateField, value: "VRT", want: InvalidWarning},
			values:       map[string]string{"MY_DXCC": "Slovak Republic", "MY_STATE": "VRT"},
		},
		// PM in Russia, but Cyrilic characters
		{
			validateTest: validateTest{field: StateField, value: "ПМ", want: InvalidError},
			values:       map[string]string{"DXCC": "15", "STATE": "ПМ"},
		},
	}
	for _, tc := range tests {
		ctx := ValidationContext{FieldValue: func(name string) string { return tc.values[name] }}
		testValidator(t, tc.validateTest, ctx, "ValidateEnumScope")
	}
}

func TestValidateStringEnumScope(t *testing.T) {
	// Submode is a String field but has an associated enumeration
	tests := []struct {
		validateTest
		values map[string]string
	}{
		{
			validateTest: validateTest{field: SubmodeField, value: "", want: Valid},
			values:       map[string]string{},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "", want: Valid},
			values:       map[string]string{"MODE": "CW", "SUBMOE": ""},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "THOR-M", want: Valid},
			values:       map[string]string{"MODE": "", "SUBMOE": "THOR-M"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "LSB", want: Valid},
			values:       map[string]string{"MODE": "SSB", "SUBMOE": "LSB"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "PSK31", want: Valid},
			values:       map[string]string{"MODE": "PSK", "SUBMOE": "PSK31"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "Dstar", want: Valid},
			values:       map[string]string{"MODE": "DIGITALVOICE", "SUBMOE": "Dstar"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "OLIVIA 16/500", want: Valid},
			values:       map[string]string{"MODE": "OLIVIA", "SUBMOE": "OLIVIA 16/500"},
		},

		{
			validateTest: validateTest{field: SubmodeField, value: "NOTAMODE", want: InvalidWarning},
			values:       map[string]string{},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "lower", want: InvalidWarning},
			values:       map[string]string{"MODE": "SSB", "SUBMOE": "lower"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "LSB", want: InvalidWarning},
			values:       map[string]string{"MODE": "FM", "SUBMOE": "LSB"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "PSK31", want: InvalidWarning},
			values:       map[string]string{"MODE": "CW", "SUBMOE": "PSK31"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "VOCODER", want: InvalidWarning},
			values:       map[string]string{"MODE": "DIGITALVOICE", "SUBMOE": "VOCODER"},
		},
	}
	for _, tc := range tests {
		ctx := ValidationContext{FieldValue: func(name string) string { return tc.values[name] }}
		testValidator(t, tc.validateTest, ctx, "ValidateStringScope")
	}
}

func TestValidateContestID(t *testing.T) {
	// CONTEST_ID is associated with the Contest_Id enum but it's a String field,
	// allowing non-enumerated contests (but encouraging enum values for
	// interoperability).
	tests := []validateTest{
		{field: ContestIdField, value: "", want: Valid},
		{field: ContestIdField, value: "CO-QSO-PARTY", want: Valid},
		{field: ContestIdField, value: "CO_QSO_PARTY", want: InvalidWarning},
		{field: ContestIdField, value: "CO QSO PARTY", want: InvalidWarning},
		// Enum value for IL has spaces rather than dashes for some reason
		{field: ContestIdField, value: "IL QSO PARTY", want: Valid},
		{field: ContestIdField, value: "IL-QSO-PARTY", want: InvalidWarning},
		{field: ContestIdField, value: "IL_QSO_PARTY", want: InvalidWarning},
		{field: ContestIdField, value: "070-31-FLAVORS", want: Valid},
		{field: ContestIdField, value: "070-31-flavors", want: Valid},
		{field: ContestIdField, value: "70-31-FLAVORS", want: InvalidWarning},
		{field: ContestIdField, value: "RAC", want: Valid},
		{field: ContestIdField, value: "R.A.C.", want: InvalidWarning},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateContestId")
	}
}
