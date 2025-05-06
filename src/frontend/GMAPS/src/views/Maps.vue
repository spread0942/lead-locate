<template>
    <h1>Maps</h1>

    <input type="text" placeholder="Search for a place" v-model="searchText" />
    <input type="text" placeholder="Search your targer" v-model="searchTarget" />
    <input type="text" placeholder="Limit" v-model="searchLimit" />
    <button @click="search">
        <span v-if="isLoading">
            <i class="fa-solid fa-spinner fa-spin"></i>
        </span>
        <span v-else>
            <i class="fa-solid fa-arrow-right"></i>
        </span>
        Search
    </button>
    <button @click="downloadXlsx" :disabled="companies.length === 0" >Download XLSX</button>

    <div style="overflow-x: auto; overflow-y: auto; height: 400px;">
        <table>
            <thead>
                <tr>
                    <th style="width: 150px;">Title</th>
                    <th style="width: 300px;">Website</th>
                    <th style="width: 100px;">Phone</th>
                    <th style="width: 300px;">Address</th>
                    <th style="width: 300px;">Place ID</th>
                    <th style="width: 100px;">Category</th>
                    <th style="width: 80px;">Rating</th>
                    <th style="width: 80px;">Gowler</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(company, index) in companies" :key="index">
                    <td>{{ company.title }}</td>
                    <td><a :href="company.website" target="_blank">{{ company.website }}</a></td>
                    <td>{{ company.phone }}</td>
                    <td>{{ company.address }}</td>
                    <td>{{ company.placeId }}</td>
                    <td>{{ company.category }}</td>
                    <td>{{ company.rating }}</td>
                    <td>
                        <button @click="gowlerTheSite(company)">
                            <i class="fa-solid fa-spider"></i>
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>

    <div style="overflow-x: auto; overflow-y: auto; height: 400px;margin-top: 20px;">
        <table>
            <thead>
                <tr>
                    <th style="width: 300px;">Coordinate</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(cordinate, index) in coordinates" :key="index">
                    <td>{{ cordinate }}</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script>
    import { ref } from 'vue';
    import { useRouter } from 'vue-router';

    import * as XLSX from 'xlsx';

    export default {
        name: 'Maps',
        components: { },
        setup() {
            const router = useRouter();
            const searchText = ref('');
            const searchTarget = ref('');
            const searchLimit = ref('10');
            const companies = ref([]);
            const coordinates = ref([]);
            const isLoading = ref(false);

            // Function to query Google Maps API
            const search = async () => {
                isLoading.value = true;
                const params = new URLSearchParams(
                    {
                        location: searchText.value,
                        target: searchTarget.value,
                        limit: searchLimit.value
                    });
                const response = await fetch(`/api/maps?${params}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });
                const data = await response.json();
                companies.value = data.companies.map(company => ({
                    title: company.title || 'N/A',
                    website: company.website || 'N/A',
                    phone: company.phone || 'N/A',
                    address: company.address || 'N/A',
                    placeId: company.placeId || 'N/A',
                    category: company.category || 'N/A',
                    rating: company.rating || 'N/A'
                }));
                coordinates.value = data.coordinates;
                isLoading.value = false;
            };

            // Function to navigate to the Gowler route
            const gowlerTheSite = (company) => {
                if (company.website) {
                    router.push({ name: 'Gowler', params: { website: company.website } });
                }
            };

            // Function to download Xlsx file
            const downloadXlsx = () => {
                const now = Date.now();
                const sheetData1 = [
                    ['Title', 'Website', 'Phone', 'Address', 'Place ID', 'Category', 'Rating'],
                    ...companies.value.map(company => [
                        company.title,
                        company.website,
                        company.phone,
                        company.address,
                        company.placeId,
                        company.category,
                        company.rating
                    ])
                ];
                const sheetData2 = [
                    ['Coordinates'],
                    ...(coordinates.value ?? []).map(coord => [coord])
                ];

                const workbook = XLSX.utils.book_new();

                const worksheet1 = XLSX.utils.aoa_to_sheet(sheetData1);
                const worksheet2 = XLSX.utils.aoa_to_sheet(sheetData2);

                XLSX.utils.book_append_sheet(workbook, worksheet1, 'Companies');
                XLSX.utils.book_append_sheet(workbook, worksheet2, 'Coordinates');

                const fileName = `companies_${now}.xlsx`;
                XLSX.writeFile(workbook, fileName);
            };

            return {
                searchText,
                searchTarget,
                searchLimit,
                companies,
                coordinates,
                isLoading,

                search,
                downloadXlsx,
                gowlerTheSite
            };
        },
    }
</script>

<style>
    i {
        font-size: 20px;
        color: white;
    }
    i:hover {
        color: #83C5BE;
    }
    a {
        color: #2c3e50;
    }
    a:hover {
        text-decoration: underline;
    }
    button i {
        font-size: 15px;
        color: white;
    }
    button i:hover {
        color: #fff;
    }
</style>