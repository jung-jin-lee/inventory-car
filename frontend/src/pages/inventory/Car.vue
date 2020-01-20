<template>
    <v-container>
        <v-row>
            <v-col cols="12">
                <v-card height="4em" outlined>
                    <v-card-title>
                        {{ pageTitle }}
                    </v-card-title>
                </v-card>
            </v-col>
        </v-row>
        <v-row>
            <v-col cols="9">
                <v-combobox
                    v-model="selectedFilterOptions"
                    :items="filterOptions"
                    label="Filter"
                    multiple
                    dense
                ></v-combobox>
            </v-col>
            <v-col cols="3">
                <v-btn>Update Filter</v-btn>
            </v-col>
        </v-row>
        <v-row>
            <v-data-table
                class="overflow-auto"
                :headers="carInventoryTableHeaders"
                hide-default-footer
                height="600"
                :items="addNoInCarInventories"
                :items-per-page="-1"
                item-key="vin"
                no-data-text="Inventory cars are empty"
                show-select
                :single-select="singleSelect"
                v-model="selectedRows"
            >
                <template v-slot:item.msrp="{ item }">
                    <span>{{ item.msrp | numberFormat }}</span>
                </template>
            </v-data-table>
        </v-row>
        <v-row>
            <v-col cols="12">
                <v-card class="pt-4 pl-4" outlined>
                    <v-row>
                        <v-col xs="5" md="3" lg="3" xl="4">
                            <v-btn class="inline-block" outlined tile @click="onClickAddRowButton">+</v-btn>
                            <v-btn class="inline-block" outlined tile @click="onClickRemoveRowButton">-</v-btn>
                        </v-col>
                        <v-col xs="5" md="3" lg="3" xl="4">
                            <v-file-input  outlineds dense tile label="Upload File"></v-file-input>
                        </v-col>
                        <v-spacer></v-spacer>
                    </v-row>
                </v-card>
            </v-col>
        </v-row>
        <v-row justify="center">
            <v-dialog v-model="visibleDialog" persistent max-width="600px">
                <v-card>
                    <v-card-title>
                        <span class="headline">Add Inventory Car Item</span>
                    </v-card-title>
                    <v-card-text>
                        <v-form ref="formInDialog" v-model="isFormValidInDialog">
                            <v-container>
                                <v-row>
                                    <v-col cols="12">
                                        <v-text-field
                                            label="VIN#"
                                            required
                                            :rules="basicRequiredRules"
                                            v-model="enteredVin"
                                        ></v-text-field>
                                    </v-col>
                                    <v-col cols="12">
                                        <v-text-field 
                                            label="Model"
                                            required
                                            :rules="basicRequiredRules"
                                            v-model="enteredModel"
                                        ></v-text-field>
                                    </v-col>
                                    <v-col cols="12">
                                        <v-text-field 
                                            label="Make"
                                            required
                                            :rules="basicRequiredRules"
                                            v-model="enteredMake"
                                        ></v-text-field>
                                    </v-col>
                                    <v-col cols="12">
                                        <v-select
                                            :items="years"
                                            label="Year"
                                            required
                                            :value="selectedYear"
                                            @input="onInputYearInDialog"
                                        ></v-select>
                                    </v-col>
                                    <v-col cols="12">
                                        <v-text-field
                                            label="MSRP"
                                            required
                                            :rules="notNegativeIntegerRules"
                                            type="number"
                                            :value="enteredMSRP"
                                            @input="$event => onInputMSRPInDialog($event)"
                                            @keydown="onKeydownMSRPInDialog"
                                        ></v-text-field>
                                    </v-col>
                                    <v-col cols="12">
                                        <v-select
                                        :items="['ordered', 'in stock', 'sold']"
                                        label="Status"
                                        required
                                        v-model="selectedStatus"
                                        ></v-select>
                                    </v-col>
                                    <v-col cols="6" sm="6">
                                        <v-select
                                        :items="['y', 'n']"
                                        label="Booked"
                                        required
                                        v-model="selectedBooked"
                                        ></v-select>
                                    </v-col>
                                    <v-col cols="12" sm="6">
                                        <v-select
                                        :items="['y', 'n']"
                                        label="Listed"
                                        required
                                        v-model="selectedListed"
                                        ></v-select>
                                    </v-col>
                                </v-row>
                            </v-container>
                        </v-form>
                    </v-card-text>
                    <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" text @click="onClickInventoryCarSaveButton">Save</v-btn>
                    <v-btn color="blue darken-1" text @click="hideDialog">Cancel</v-btn>
                    </v-card-actions>
                </v-card>
            </v-dialog>
        </v-row>
        <v-snackbar
            bottom
            :color="toastMessageLevel"
            :multi-line="multiLineToastMessage"
            :timeout="timeout"
            v-model="showToastMessage"
        >
            {{ toastMessage }}
            <v-btn
                color="white"
                text
                @click="showToastMessage = false"
            >
                X
            </v-btn>
        </v-snackbar>
    </v-container>
</template>

<script>
import { getInventoryCars, createInventoryCar, deleteInventoryCar } from "../../api/inventory/car";
import { handleApiResponse } from "../../utils/common/api-handler";

export default {
    computed: {
        addNoInCarInventories() {
            return this.carInventories.length > 0 
                ?  this.carInventories.map((v, i) => ({ ...v, no: i + 1 }))
                :  this.carInventories;
        },
    },
    data() {
        const generateYearData = (startYear, endYear = startYear) => {
            if (endYear < startYear) {
                throw new Error(`endYear must bigger than startYear. { "startYear": ${startYear}, "endYear": ${endYear} }`);
            }
            const n = endYear - startYear + 1;
            return [...new Array(n)].map(() => ({ "text": String(startYear), "value": startYear++ }));
        };
        const notEmptyRule = v => !!v || "This field is required";
        const notNegativeIntegerRule = v => (Number(v) >= 0 || "This field is non negative integer.");

        const DEFAULT_STATUS = "ordered";
        const DEFAULT_BOOKED = "y";
        const DEFAULT_LISTED = "y";
        const START_YEAR = 1970;
        const END_YEAR = new Date().getFullYear();
        const yearData = generateYearData(START_YEAR, END_YEAR);

        return ({
            basicRequiredRules: [notEmptyRule],
            carInventories: [],
            carInventoryTableHeaders: [
                { "align": "left",   "text": "No",     "value": "no",     "width": 80  },
                { "align": "left",   "text": "Vin#",   "value": "vin",    "width": 200 },
                { "align": "left",   "text": "Model",  "value": "model",  "width": 140 },
                { "align": "left",   "text": "Make",   "value": "make",   "width": 100 },
                { "align": "center", "text": "Year",   "value": "year",   "width": 80  },
                { "align": "right",  "text": "MSRP",   "value": "msrp",   "width": 100 },
                { "align": "left",   "text": "Status", "value": "status", "width": 140 },
                { "align": "center", "text": "Booked", "value": "booked", "width": 100 },
                { "align": "center", "text": "Listed", "value": "listed", "width": 100 },
            ],
            visibleDialog: false,
            disablePagination: true,
            enteredMake: "",
            enteredModel: "",
            enteredMSRP: 0,
            enteredVin: "",
            enteredYear: 0,
            fields: [
                { "prop": "enteredVin",     "name": "vin",    "defaultValue": "" },
                { "prop": "enteredModel",   "name": "model",  "defaultValue": "" },
                { "prop": "enteredMake",    "name": "make",   "defaultValue": "" },
                { "prop": "selectedYear",   "name": "year",   "defaultValue": END_YEAR },
                { "prop": "enteredMSRP",    "name": "msrp",   "defaultValue": 0 },
                { "prop": "selectedStatus", "name": "status", "defaultValue": DEFAULT_STATUS },
                { "prop": "selectedBooked", "name": "booked", "defaultValue": DEFAULT_BOOKED },
                { "prop": "selectedListed", "name": "listed", "defaultValue": DEFAULT_LISTED },
            ],
            filterOptions: [
                { "text": "Model", "value": "model" },
                { "text": "Make",  "value": "make"  },
                { "text": "Year",  "value": "year"  },
            ],
            isFormValidInDialog: false,
            logger: console,
            multiLineToastMessage: false,
            notNegativeIntegerRules: [notNegativeIntegerRule],
            pageTitle: "Inventory List",
            selectedBooked: DEFAULT_BOOKED,
            selectedFilterOptions: [],
            selectedListed: DEFAULT_LISTED,
            selectedRows: [],
            selectedStatus: DEFAULT_STATUS,
            selectedYear: END_YEAR,
            showToastMessage: false,
            singleSelect: true,
            timeout: 2000,
            toastMessage: "",
            toastMessageLevel: "success",
            years: yearData,
        })
    },
    filters: {
        numberFormat(number) {
            if (!number) {
                return "0";
            }
            
            return new Intl.NumberFormat("en-US").format(number);
        },
    },
    methods: {
        getInventoryCarDataFromDialog(fields) {
            return fields.reduce((data, f) => {
                data[f.name] = this[f.prop];

                return data;
            }, {});
        },
        hideDialog() {
            this.visibleDialog = false;
        },
        initFieldsInDialog(fields) {
            fields.forEach(f => this[f.prop] = f.defaultValue);
        },
        isNotNumber(value) {
            return !/^[0-9]*$/gm.test(value);
        },
        onClickAddRowButton() {
            this.showDialog();
        },
        async onClickInventoryCarSaveButton() {
            try {
                if (!this.$refs.formInDialog.validate()) {
                    this.openToastMessage("some fields are invalid.", "red");
                    return false;
                }
                
                await createInventoryCar(
                    this.getInventoryCarDataFromDialog(this.fields)
                );
                this.carInventories = await getInventoryCars().then(handleApiResponse);
                this.openToastMessage("Added Row");
                this.hideDialog();
                this.initFieldsInDialog(this.fields);
                this.$refs.formInDialog.resetValidation();
            } catch (error) {
                this.openToastMessageAboutServerError();
                this.logger.error("[API][createInventoryCar or getInventoryCars]error:", error);
            }
        },
        async onClickRemoveRowButton() {
            if (this.selectedRows.length > 0) {
                const [{ vin }] = this.selectedRows;
                try {
                    await deleteInventoryCar(vin);
                    this.carInventories = await getInventoryCars().then(handleApiResponse);
                    this.openToastMessage("Deleted Checked Row");
                } catch (error) {
                    this.openToastMessageAboutServerError();
                    this.logger.error("[API][deleteInventoryCar or getInventoryCars]", error);
                }
            } else {
                this.openToastMessage("check row!", "red");
            }      
        },
        onInputMSRPInDialog(enteredMSRP) {
            this.enteredMSRP = Number(enteredMSRP);
        },
        onInputYearInDialog(selectedYear) {
            this.selectedYear = Number(selectedYear);
        },
        onKeydownMSRPInDialog($event) {
            if ($event.ctrlKey && $event.key === "a") {
                return true;
            }

            if ($event.key === "Home" || $event.key === "End") {
                return true;
            }

            if ($event.key === "Backspace") {
                return true;
            }
            
            if (this.isNotNumber($event.key)) {
                $event.preventDefault();
                return false;
            }
        },
        openToastMessage(message, level = "success", multiLine = false) {
            this.toastMessage = message;
            this.toastMessageLevel = level;
            this.multiLine = multiLine;
            this.showToastMessage = true;
        },
        openToastMessageAboutServerError() {
            this.openToastMessage(
                `There was a problem communicating with the server.
                 Please Check your network environment.`, 
                "red",
                true,
            );
        },
        showDialog() {
            this.visibleDialog = true;
        }
    },
    async mounted() {
        try {
            this.carInventories = await getInventoryCars().then(handleApiResponse);
        } catch (error) {
            this.openToastMessageAboutServerError();
            this.logger.error("[API][getInventoryCars]", error);
        }
    },
}
</script>

<style>
.overflow-auto {
    overflow: auto;
}
.inline-block {
    display: inline-block;
}
.pt-4 {
    padding-top: 1em;
}
.pl-4 {
    padding-left: 1em;
}
</style>