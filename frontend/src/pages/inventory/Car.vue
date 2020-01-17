<template>
    <v-container>
        <v-row>
            <v-col cols="12">
                <v-card height="4em" outlined>
                    <v-card-title>
                        {{pageTitle}}
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
                :items="carInventories"
                :items-per-page="-1"
                show-select
            ></v-data-table>
        </v-row>
        <v-row>
            <v-col xs="6" md="3">
                <v-btn style="display: inline-block;" outlined tile @click="onClickAddRowButton">+</v-btn>
                <v-btn style="display: inline-block;"  outlined tile>-</v-btn>
            </v-col>
            <v-col xs="6" md="3">
                <v-file-input outlined dense tile label="Upload File"></v-file-input>
            </v-col>
            <v-spacer></v-spacer>
        </v-row>
        <v-row justify="center">
            <v-dialog v-model="dialog" persistent max-width="600px">
                <v-card>
                    <v-card-title>
                        <span class="headline">Add Car Inventory Item</span>
                    </v-card-title>
                    <v-card-text>
                        <v-container>
                            <v-row>
                                <v-col cols="12">
                                    <v-text-field label="VIN#" required></v-text-field>
                                </v-col>
                                <v-col cols="12">
                                    <v-text-field label="Model" required></v-text-field>
                                </v-col>
                                <v-col cols="12">
                                    <v-text-field label="Make" required></v-text-field>
                                </v-col>
                                <v-col cols="12">
                                    <v-text-field label="Year" required></v-text-field>
                                </v-col>
                                <v-col cols="12">
                                    <v-text-field label="MSRP" required></v-text-field>
                                </v-col>
                                <v-col cols="6" sm="6">
                                    <v-select
                                    :items="['y', 'n']"
                                    label="Booked"
                                    required
                                    ></v-select>
                                </v-col>
                                <v-col cols="12" sm="6">
                                    <v-select
                                    :items="['y', 'n']"
                                    label="Listed"
                                    required
                                    ></v-select>
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-card-text>
                    <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" text @click="dialog = false">Save</v-btn>
                    <v-btn color="blue darken-1" text @click="dialog = false">Cancel</v-btn>
                    </v-card-actions>
                </v-card>
            </v-dialog>
        </v-row>
    </v-container>
</template>

<script>
export default {
    data: () => {
        const sampleData = {
            no: 1,
            vinNo: "MEFLWFSDFKMLELKF",
            model: "320i",
            make: "BMW",
            year: 2014,
            msrp: 147000,
            status: "ordered",
            booked: "y",
            listed: "n",
        };
        const generateSampleData = (sampleData, n) => ([...new Array(n)].map(() => sampleData));

        return ({
        carInventories: generateSampleData(sampleData, 20),
        carInventoryTableHeaders: [
            { align: "left", text: "No", value: "no", width: 40 },
            { align: "left", text: "Vin#", value: "vinNo", width: 200 },
            { align: "left", text: "Model", value: "model", width: 140 },
            { align: "left", text: "Make", value: "make", width: 100 },
            { align: "center", text: "Year", value: "year", width: 80 },
            { align: "right", text: "MSRP", value: "msrp", width: 100 },
            { align: "left", text: "Status", value: "status", width: 140 },
            { align: "center", text: "Booked", value: "booked", width: 80 },
            { align: "center", text: "Listed", value: "listed", width: 80 },
        ],
        dialog: false,
        disablePagination: true,
        filterOptions: [
            {
                "text": "Model",
                "value": "model"
            },
            {
                "text": "Make",
                "value": "make"
            },
            {
                "text": "Year",
                "value": "year"
            },
        ],
        pageTitle: "Inventory List",
        selectedFilterOptions: [],
    })},
    methods: {
        onClickAddRowButton() {
            this.dialog = true;
        }
    }
}
</script>

<style>
.overflow-auto {
    overflow: auto;
}
</style>